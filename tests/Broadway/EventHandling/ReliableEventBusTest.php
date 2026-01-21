<?php

declare(strict_types=1);

namespace Vonq\Webshop\Tests\Broadway\EventHandling;

use Prophecy\PhpUnit\ProphecyTrait;
use Broadway\Domain\DomainEventStream;
use Broadway\Domain\DomainMessage;
use Broadway\Domain\Metadata;
use Broadway\EventHandling\EventListener;
use Broadway\EventHandling\ReliableEventBus;
use Monolog\Handler\TestHandler;
use Monolog\Logger;
use PHPUnit\Framework\MockObject\MockObject;
use PHPUnit\Framework\TestCase;
use Psr\Log\LoggerInterface;

class ReliableEventBusTest extends TestCase
{
    use ProphecyTrait;

    /**
     * @var ReliableEventBus
     */
    private $eventBus;

    /** @var LoggerInterface */
    private $logger;

    /** @var TestHandler */
    private $testHandler;

    protected function setUp(): void
    {
        $this->logger = new Logger('main');
        $this->testHandler = new TestHandler();
        $this->logger->pushHandler($this->testHandler);

        $this->eventBus = new ReliableEventBus($this->logger);
    }

    /**
     * @test
     */
    public function it_should_process_the_next_handle_when_a_handler_fails(): void
    {
        $domainMessage = $this->createDomainMessage([]);
        $domainEventStream = new DomainEventStream([$domainMessage]);

        $eventListener1 = $this->createMock(EventListener::class);
        $eventListener1
            ->expects($this->once())
            ->method('handle')
            ->with($domainMessage)
            ->willThrowException(new \Exception());

        $eventListener2 = $this->createMock(EventListener::class);
        $eventListener2
            ->expects($this->once())
            ->method('handle')
            ->with($domainMessage);

        $this->eventBus->subscribe($eventListener1);
        $this->eventBus->subscribe($eventListener2);

        $this->eventBus->publish($domainEventStream);

        $this->assertTrue(
            $this->testHandler->hasErrorThatContains(
                sprintf('[Event LISTENER]: %s', $eventListener1::class)
            )
        );
    }


    /**
     * @test
     */
    public function it_subscribes_an_event_listener()
    {
        $domainMessage = $this->createDomainMessage(['foo' => 'bar']);

        $eventListener = $this->createEventListenerMock();
        $eventListener
            ->expects($this->once())
            ->method('handle')
            ->with($domainMessage);

        $this->eventBus->subscribe($eventListener);

        $this->eventBus->publish(new DomainEventStream([$domainMessage]));
    }

    /**
     * @test
     */
    public function it_publishes_events_to_subscribed_event_listeners(): void
    {
        $domainMessage1 = $this->createDomainMessage([]);
        $domainMessage2 = $this->createDomainMessage([]);

        $domainEventStream = new DomainEventStream([$domainMessage1, $domainMessage2]);

        $calls1 = [];
        $eventListener1 = $this->createEventListenerMock();
        $eventListener1
            ->expects($this->exactly(2))
            ->method('handle')
            ->willReturnCallback(function ($message) use (&$calls1): void {
                $calls1[] = $message;
            });

        $calls2 = [];
        $eventListener2 = $this->createEventListenerMock();
        $eventListener2
            ->expects($this->exactly(2))
            ->method('handle')
            ->willReturnCallback(function ($message) use (&$calls2): void {
                $calls2[] = $message;
            });

        $this->eventBus->subscribe($eventListener1);
        $this->eventBus->subscribe($eventListener2);

        $this->eventBus->publish($domainEventStream);

        $this->assertSame([$domainMessage1, $domainMessage2], $calls1);
        $this->assertSame([$domainMessage1, $domainMessage2], $calls2);
    }


    /**
     * @test
     */
    public function it_does_not_dispatch_new_events_before_all_listeners_have_run(): void
    {
        $domainMessage1 = $this->createDomainMessage(['foo' => 'bar']);
        $domainMessage2 = $this->createDomainMessage(['foo' => 'bas']);

        $domainEventStream = new DomainEventStream([$domainMessage1]);

        $eventListener1 = new SimpleEventBusTestListener(
            $this->eventBus,
            new DomainEventStream([$domainMessage2])
        );

        $calls = [];
        $eventListener2 = $this->createEventListenerMock();
        $eventListener2
            ->expects($this->exactly(2))
            ->method('handle')
            ->willReturnCallback(function ($message) use (&$calls): void {
                $calls[] = $message;
            });

        $this->eventBus->subscribe($eventListener1);
        $this->eventBus->subscribe($eventListener2);

        $this->eventBus->publish($domainEventStream);

        // Proprietà chiave: prima gestisce domainMessage1, poi (solo dopo) domainMessage2
        $this->assertSame([$domainMessage1, $domainMessage2], $calls);
    }

    /**
     * @test
     */
    public function it_should_still_publish_events_after_exception(): void
    {
        $domainMessage1 = $this->createDomainMessage(['foo' => 'bar']);
        $domainMessage2 = $this->createDomainMessage(['foo' => 'bas']);

        $domainEventStream1 = new DomainEventStream([$domainMessage1]);
        $domainEventStream2 = new DomainEventStream([$domainMessage2]);

        $calls = [];
        $eventListener = $this->createEventListenerMock();
        $eventListener
            ->expects($this->exactly(2))
            ->method('handle')
            ->willReturnCallback(function ($message) use (&$calls, $domainMessage1): void {
                $calls[] = $message;

                if ($message == $domainMessage1) {
                    throw new \Exception('I failed.');
                }
            });

        $this->eventBus->subscribe($eventListener);

        // ReliableEventBus tipicamente "swallow-a" l'eccezione e continua (loggando).
        // Quindi NON facciamo try/catch: se rilanciasse, il test fallirebbe qui.
        $this->eventBus->publish($domainEventStream1);
        $this->eventBus->publish($domainEventStream2);

        $this->assertSame([$domainMessage1, $domainMessage2], $calls);
    }

    private function createEventListenerMock(): MockObject
    {
        return $this->createMock(EventListener::class);
    }

    private function createDomainMessage($payload)
    {
        return DomainMessage::recordNow(1, 1, new Metadata([]), new SimpleEventBusTestEvent($payload));
    }
}

class SimpleEventBusTestEvent
{
    public $data;

    public function __construct($data)
    {
        $this->data = $data;
    }
}

class SimpleEventBusTestListener implements EventListener
{
    private $eventBus;
    private $handled = false;
    private $publishableStream;

    public function __construct($eventBus, $publishableStream)
    {
        $this->eventBus = $eventBus;
        $this->publishableStream = $publishableStream;
    }

    public function handle(DomainMessage $domainMessage): void
    {
        if (!$this->handled) {
            $this->eventBus->publish($this->publishableStream);
            $this->handled = true;
        }
    }
}

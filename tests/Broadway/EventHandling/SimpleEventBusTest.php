<?php

/*
 * This file is part of the broadway/broadway package.
 *
 * (c) 2020 Broadway project
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Broadway\EventHandling;

use Broadway\Domain\DomainEventStream;
use Broadway\Domain\DomainMessage;
use Broadway\Domain\Metadata;
use PHPUnit\Framework\MockObject\MockObject;
use PHPUnit\Framework\TestCase;

class SimpleEventBusTest extends TestCase
{
    private SimpleEventBus $eventBus;

    protected function setUp(): void
    {
        $this->eventBus = new SimpleEventBus();
    }

    /** @test */
    public function it_subscribes_an_event_listener(): void
    {
        $domainMessage = $this->createDomainMessage(['foo' => 'bar']);

        $eventListener = $this->createEventListenerMock();
        $eventListener
            ->expects($this->once())
            ->method('handle')
            ->with($this->domainMessageMatchesPayload(['foo' => 'bar']));

        $this->eventBus->subscribe($eventListener);
        $this->eventBus->publish(new DomainEventStream([$domainMessage]));
    }

    /** @test */
    public function it_publishes_events_to_subscribed_event_listeners(): void
    {
        $domainMessage1 = $this->createDomainMessage(['n' => 1]);
        $domainMessage2 = $this->createDomainMessage(['n' => 2]);

        $domainEventStream = new DomainEventStream([$domainMessage1, $domainMessage2]);

        $calls1 = [];
        $eventListener1 = $this->createEventListenerMock();
        $eventListener1
            ->expects($this->exactly(2))
            ->method('handle')
            ->willReturnCallback(function (DomainMessage $message) use (&$calls1): void {
                $calls1[] = $message;
            });

        $calls2 = [];
        $eventListener2 = $this->createEventListenerMock();
        $eventListener2
            ->expects($this->exactly(2))
            ->method('handle')
            ->willReturnCallback(function (DomainMessage $message) use (&$calls2): void {
                $calls2[] = $message;
            });

        $this->eventBus->subscribe($eventListener1);
        $this->eventBus->subscribe($eventListener2);

        $this->eventBus->publish($domainEventStream);

        $this->assertCount(2, $calls1);
        $this->assertCount(2, $calls2);

        $this->assertSame(['n' => 1], $calls1[0]->getPayload()->data);
        $this->assertSame(['n' => 2], $calls1[1]->getPayload()->data);

        $this->assertSame(['n' => 1], $calls2[0]->getPayload()->data);
        $this->assertSame(['n' => 2], $calls2[1]->getPayload()->data);
    }

    /** @test */
    public function it_does_not_dispatch_new_events_before_all_listeners_have_run(): void
    {
        $domainMessage1 = $this->createDomainMessage(['foo' => 'bar']);
        $domainMessage2 = $this->createDomainMessage(['foo' => 'bas']);

        // Pubblico solo il primo messaggio.
        $domainEventStream = new DomainEventStream([$domainMessage1]);

        // Listener1 durante la prima handle() pubblica un nuovo stream contenente domainMessage2.
        $eventListener1 = new SimpleEventBusTestListener(
            $this->eventBus,
            new DomainEventStream([$domainMessage2])
        );

        // Listener2 deve ricevere prima message1 e solo dopo message2.
        $calls = [];
        $eventListener2 = $this->createEventListenerMock();
        $eventListener2
            ->expects($this->exactly(2))
            ->method('handle')
            ->willReturnCallback(function (DomainMessage $message) use (&$calls): void {
                $calls[] = $message;
            });

        $this->eventBus->subscribe($eventListener1);
        $this->eventBus->subscribe($eventListener2);

        $this->eventBus->publish($domainEventStream);

        $this->assertCount(2, $calls);

        // Se il payload è SimpleEventBusTestEvent con proprietà ->data:
        $this->assertSame(['foo' => 'bar'], $calls[0]->getPayload()->data);
        $this->assertSame(['foo' => 'bas'], $calls[1]->getPayload()->data);

        // In alternativa, se getPayload() restituisce direttamente l'array, usa:
        // $this->assertSame(['foo' => 'bar'], $calls[0]->getPayload());
        // $this->assertSame(['foo' => 'bas'], $calls[1]->getPayload());
    }

    /** @test */
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
            ->willReturnCallback(function (DomainMessage $message) use (&$calls): void {
                $calls[] = $message;

                // Lancia solo sul primo payload
                if ($message->getPayload()->data === ['foo' => 'bar']) {
                    throw new \Exception('I failed.');
                }
            });

        $this->eventBus->subscribe($eventListener);

        try {
            $this->eventBus->publish($domainEventStream1);
            $this->fail('Expected exception was not thrown.');
        } catch (\Exception $e) {
            $this->assertSame('I failed.', $e->getMessage());
        }

        $eventListener2 = $this->createEventListenerMock();
        $eventListener2
            ->expects($this->once())
            ->method('handle')
            ->willReturnCallback(function (DomainMessage $message): void {
                $this->assertSame(['foo' => 'bas'], $message->getPayload()->data);
            });

        $this->eventBus->subscribe($eventListener2);

        // Nonostante l'eccezione prima, questo publish deve completare e chiamare i listener.
        $this->eventBus->publish($domainEventStream2);

        // Verifichiamo anche l'ordine delle chiamate del primo listener
        $this->assertCount(2, $calls);
        $this->assertSame(['foo' => 'bar'], $calls[0]->getPayload()->data);
        $this->assertSame(['foo' => 'bas'], $calls[1]->getPayload()->data);
    }


    private function createEventListenerMock(): MockObject
    {
        return $this->createMock(EventListener::class);
    }

    private function createDomainMessage(array $payload): DomainMessage
    {
        // recordedOn è "now": non lo useremo mai per confronti stretti nei test.
        return DomainMessage::recordNow(
            1,
            1,
            new Metadata([]),
            new SimpleEventBusTestEvent($payload)
        );
    }

    /**
     * Matcher: verifica il payload del DomainMessage senza confrontare recordedOn (microsecondi).
     */
    private function domainMessageMatchesPayload(array $expectedPayload): \PHPUnit\Framework\Constraint\Constraint
    {
        return $this->callback(function ($arg) use ($expectedPayload): bool {
            if (!$arg instanceof DomainMessage) {
                return false;
            }

            $payload = $arg->getPayload();
            if (!$payload instanceof SimpleEventBusTestEvent) {
                return false;
            }

            return $payload->data === $expectedPayload;
        });
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

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

namespace Broadway\CommandHandling;

use PHPUnit\Framework\MockObject\MockObject;
use PHPUnit\Framework\TestCase;

class SimpleCommandBusTest extends TestCase
{
    /**
     * @var SimpleCommandBus
     */
    private $commandBus;

    protected function setUp(): void
    {
        $this->commandBus = new SimpleCommandBus();
    }

    /**
     * @test
     */
    public function it_dispatches_commands_to_subscribed_handlers()
    {
        $command = ['Hi' => 'There'];

        $this->commandBus->subscribe($this->createCommandHandlerMock($command));
        $this->commandBus->subscribe($this->createCommandHandlerMock($command));
        $this->commandBus->dispatch($command);
    }

    /**
     * @test
     */
    public function it_does_not_handle_new_commands_before_all_commandhandlers_have_run(): void
    {
        $command1 = ['foo' => 'bar'];
        $command2 = ['foo' => 'bas'];

        $calls = [];

        $commandHandler = $this->createMock(CommandHandler::class);
        $commandHandler
            ->expects($this->exactly(2))
            ->method('handle')
            ->willReturnCallback(function (array $command) use (&$calls): void {
                $calls[] = $command;
            });

        $this->commandBus->subscribe(new SimpleCommandBusTestHandler($this->commandBus, $command2));
        $this->commandBus->subscribe($commandHandler);

        $this->commandBus->dispatch($command1);

        $this->assertSame([$command1, $command2], $calls);
    }

    /**
     * @test
     */
    public function it_should_still_handle_commands_after_exception(): void
    {
        $command1 = ['foo' => 'bar'];
        $command2 = ['foo' => 'bas'];

        $failingHandler = $this->createMock(CommandHandler::class);
        $simpleHandler  = $this->createMock(CommandHandler::class);

        // Il failing handler potrebbe lanciare su command1 (ma non assumiamo altro)
        $failingHandler
            ->method('handle')
            ->willReturnCallback(function (array $command) use ($command1): void {
                if ($command === $command1) {
                    throw new \Exception('I failed.');
                }
            });

        // Ci interessa che command2 venga gestito dopo l'eccezione
        $simpleHandler
            ->expects($this->once())
            ->method('handle')
            ->with($command2);

        $this->commandBus->subscribe($failingHandler);
        $this->commandBus->subscribe($simpleHandler);

        try {
            $this->commandBus->dispatch($command1);
            $this->fail('Expected exception was not thrown.');
        } catch (\Exception $e) {
            $this->assertSame('I failed.', $e->getMessage());
        }

        $this->commandBus->dispatch($command2);
    }


    private function createCommandHandlerMock(array $expectedCommand): MockObject
    {
        $mock = $this->createMock(CommandHandler::class);

        $mock
            ->expects($this->once())
            ->method('handle')
            ->with($expectedCommand);

        return $mock;
    }
}

class SimpleCommandBusTestHandler implements CommandHandler
{
    private $commandBus;
    private $handled = false;
    private $dispatchableCommand;

    public function __construct($commandBus, $dispatchableCommand)
    {
        $this->commandBus = $commandBus;
        $this->dispatchableCommand = $dispatchableCommand;
    }

    public function handle($command): void
    {
        if (!$this->handled) {
            $this->commandBus->dispatch($this->dispatchableCommand);
            $this->handled = true;
        }
    }
}

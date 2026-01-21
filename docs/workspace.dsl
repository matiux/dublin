workspace "Dublin" "Dublin is a PHP library providing infrastructure and testing helpers for building CQRS and Event Sourced applications (a maintained fork of Broadway)." {

    !docs documentation
    !adrs decisions

    !identifiers hierarchical

    model {
        developer = person "Developer" "Builds a CQRS/ES application using Dublin."

        phpRuntime = softwareSystem "PHP Runtime" "Runtime environment (PHP 8.4/8.5)." "External"
        database = softwareSystem "Database" "Stores event streams and read models (implementation-specific)." "External"
        messageBroker = softwareSystem "Message Broker" "Optional async transport for events." "External"

        app = softwareSystem "CQRS/ES Application" "Domain application using Dublin to implement CQRS and Event Sourcing." {
          webApi = container "Application Code" "Your application codebase (domain + infrastructure) using Dublin." "PHP"

          webApi -> phpRuntime "Runs on"
          webApi -> database "Persists/loads events & read models"
          webApi -> messageBroker "Publishes/subscribes (optional)"
        }

        dublin = softwareSystem "Dublin" "CQRS + Event Sourcing infrastructure library for PHP." "Library" {

          commandBus = container "Command Bus" "Dispatches commands to handlers (sync)." "PHP"
          eventStore = container "Event Store" "Stores and loads domain event streams." "PHP"
          eventBus = container "Event Bus" "Publishes domain events to listeners." "PHP"
          processor = container "Processor" "Coordinates consuming events and updating projections/read models." "PHP"
          readModel = container "Read Model" "Repositories and in-memory implementations for projections." "PHP"
          serializer = container "Serializer" "Serializes domain messages/payloads and metadata." "PHP"
          upcasting = container "Upcasting" "Upcasters and chains to evolve event schemas." "PHP"
          testing = container "Testing Helpers" "Scenario, test base classes and utilities for CQRS/ES." "PHP"
          uuid = container "UUID Generator" "UUID v4 generator utilities." "PHP"
          auditing = container "Auditing" "Command logging / command serialization helpers." "PHP"

          eventStore -> serializer "Serializes/deserializes stored messages"
          upcasting -> eventStore "Decorates loading (upcast streams)"
          processor -> eventBus "Subscribes/consumes events"
          processor -> readModel "Updates projections/read models"
          auditing -> commandBus "Logs/serializes commands (optional)"
    }

    developer -> app "Develops and tests"
    app -> dublin "Uses as infrastructure toolkit"

    app.webApi -> dublin.commandBus "Dispatches commands"
    app.webApi -> dublin.eventStore "Appends/loads event streams"
    app.webApi -> dublin.eventBus "Publishes events"
    app.webApi -> dublin.processor "Runs processors/projectors"
    app.webApi -> dublin.readModel "Stores/queries read models"
    app.webApi -> dublin.serializer "Serializes payloads/messages"
    app.webApi -> dublin.upcasting "Upcasts legacy events"
    app.webApi -> dublin.testing "Uses in tests"
    app.webApi -> dublin.uuid "Generates UUIDs (optional)"
    app.webApi -> dublin.auditing "Command audit (optional)"

    dublin.eventStore -> database "Persists event streams" "Implementation-specific"
    dublin.readModel -> database "Persists read models" "Implementation-specific"
    dublin.eventBus -> messageBroker "Can be adapted to async transport" "Optional"
  }

  views {
    systemContext dublin "dublin-system-context" {
      include developer
      include app
      include dublin
      include phpRuntime
      include database
      include messageBroker
      autolayout lr
      title "System Context - Dublin"
    }

    container dublin "dublin-modules" {
      include *
      autolayout lr
      title "Modules (Containers) - Dublin"
    }

    styles {
      element Person {
        shape Person
      }

      element External {
        background "#666666"
        color "#ffffff"
      }

      element Library {
        background "#1f6feb"
        color "#ffffff"
      }
    }

    theme default
  }
}

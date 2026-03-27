Feature: Update Asset Attributes

  Background:
    Given I am authenticated as a FacilityManager for Tenant "dev-tenant"
    And an Asset "Rooftop HVAC Unit" with ID "asset-001" exists for Tenant "dev-tenant"

  Scenario: Successfully update Asset attributes
    When I update Asset "asset-001" with the following details:
      | field         | value                  |
      | name          | Rooftop HVAC Unit - R1 |
      | serial_number | SN-2024-9999           |
    Then I receive a 200 OK response
    And the Asset "asset-001" has name "Rooftop HVAC Unit - R1"
    And an "AssetAttributesUpdated" domain event is published for Tenant "dev-tenant"
    And the event payload includes the changed field "name"

  Scenario: Update Asset fails when name is blank
    When I update Asset "asset-001" with the following details:
      | field | value |
      | name  |       |
    Then I receive a 422 Unprocessable Entity response
    And the response contains a validation error for field "name"

  Scenario: Update Asset returns 404 when Asset does not exist
    When I update Asset "asset-does-not-exist" with the following details:
      | field | value    |
      | name  | New Name |
    Then I receive a 404 Not Found response

  Scenario: Technician cannot update an Asset
    Given I am authenticated as a Technician for Tenant "dev-tenant"
    When I update Asset "asset-001" with the following details:
      | field | value    |
      | name  | New Name |
    Then I receive a 403 Forbidden response

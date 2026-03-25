Feature: Register an Asset

  Background:
    Given I am authenticated as a FacilityManager for Tenant "dev-tenant"

  Scenario: Successfully register an Asset
    When I register an Asset with the following details:
      | field            | value              |
      | name             | Rooftop HVAC Unit  |
      | asset_type       | HVAC               |
      | facility_id      | facility-001       |
      | serial_number    | SN-2024-0042       |
    Then the Asset is created with Canonical Status "Active"
    And an "AssetRegistered" domain event is published for Tenant "dev-tenant"
    And the response includes the new Asset ID

  Scenario: List Assets returns only Assets belonging to the authenticated Tenant
    Given an Asset "Rooftop HVAC Unit" exists for Tenant "dev-tenant"
    And an Asset "Boiler Room Unit" exists for Tenant "other-tenant"
    When I request the list of Assets
    Then the response contains "Rooftop HVAC Unit"
    And the response does not contain "Boiler Room Unit"

  Scenario: Register Asset fails when required fields are missing
    When I register an Asset with the following details:
      | field       | value |
      | name        |       |
      | asset_type  | HVAC  |
    Then I receive a 422 Unprocessable Entity response
    And the response contains a validation error for field "name"

  Scenario: Technician cannot register an Asset
    Given I am authenticated as a Technician for Tenant "dev-tenant"
    When I register an Asset with the following details:
      | field       | value             |
      | name        | Rooftop HVAC Unit |
      | asset_type  | HVAC              |
      | facility_id | facility-001      |
    Then I receive a 403 Forbidden response

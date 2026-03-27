Feature: Get Asset by ID

  Background:
    Given I am authenticated as a FacilityManager for Tenant "dev-tenant"
    And an Asset "Rooftop HVAC Unit" with ID "asset-001" exists for Tenant "dev-tenant"

  Scenario: Successfully get an Asset by ID
    When I request Asset "asset-001"
    Then I receive a 200 OK response
    And the response contains the Asset with name "Rooftop HVAC Unit"
    And the response includes the canonical status "Active"

  Scenario: Get Asset returns 404 when Asset does not exist
    When I request Asset "asset-does-not-exist"
    Then I receive a 404 Not Found response

  Scenario: Get Asset returns 404 when Asset belongs to a different Tenant
    Given an Asset "Other Tenant Asset" with ID "asset-other" exists for Tenant "other-tenant"
    When I request Asset "asset-other"
    Then I receive a 404 Not Found response

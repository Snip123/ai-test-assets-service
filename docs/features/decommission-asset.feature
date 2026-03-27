Feature: Decommission Asset

  Background:
    Given I am authenticated as a FacilityManager for Tenant "dev-tenant"
    And an Asset "Rooftop HVAC Unit" with ID "asset-001" exists for Tenant "dev-tenant" with Canonical Status "Active"

  Scenario: Successfully decommission an Active Asset
    When I decommission Asset "asset-001" with reason "End of service life"
    Then I receive a 200 OK response
    And the Asset "asset-001" has Canonical Status "Decommissioned"
    And an "AssetDecommissioned" domain event is published for Tenant "dev-tenant"
    And the event payload includes:
      | field   | value               |
      | assetId | asset-001           |
      | reason  | End of service life |

  Scenario: Cannot decommission an already-Decommissioned Asset
    Given Asset "asset-001" has Canonical Status "Decommissioned"
    When I decommission Asset "asset-001" with reason "Duplicate decommission"
    Then I receive a 409 Conflict response

  Scenario: Decommission returns 404 when Asset does not exist
    When I decommission Asset "asset-does-not-exist" with reason "Not found"
    Then I receive a 404 Not Found response

  Scenario: Technician cannot decommission an Asset
    Given I am authenticated as a Technician for Tenant "dev-tenant"
    When I decommission Asset "asset-001" with reason "Technician attempt"
    Then I receive a 403 Forbidden response

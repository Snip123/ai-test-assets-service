Feature: Set Asset Location

  Background:
    Given I am authenticated as a FacilityManager for Tenant "dev-tenant"
    And an Asset "Rooftop HVAC Unit" with ID "asset-001" exists for Tenant "dev-tenant"
    And a Facility "facility-001" exists for Tenant "dev-tenant"

  Scenario: Successfully set the Location of an Asset
    When I set the Location of Asset "asset-001" to:
      | field       | value        |
      | facility_id | facility-001 |
      | location_id | roof-level-3 |
    Then I receive a 200 OK response
    And the Asset "asset-001" has facility_id "facility-001" and location_id "roof-level-3"
    And an "AssetLocationSet" domain event is published for Tenant "dev-tenant"
    And the event payload includes:
      | field      | value        |
      | assetId    | asset-001    |
      | facilityId | facility-001 |
      | locationId | roof-level-3 |

  Scenario: Asset Location can be updated to a new Location
    Given Asset "asset-001" has location_id "roof-level-3" in Facility "facility-001"
    When I set the Location of Asset "asset-001" to:
      | field       | value        |
      | facility_id | facility-001 |
      | location_id | roof-level-4 |
    Then I receive a 200 OK response
    And the Asset "asset-001" has location_id "roof-level-4"
    And an "AssetLocationSet" domain event is published for Tenant "dev-tenant"

  Scenario: Set Asset Location returns 404 when Asset does not exist
    When I set the Location of Asset "asset-does-not-exist" to:
      | field       | value        |
      | facility_id | facility-001 |
      | location_id | roof-level-3 |
    Then I receive a 404 Not Found response

  Scenario: Technician cannot set Asset Location
    Given I am authenticated as a Technician for Tenant "dev-tenant"
    When I set the Location of Asset "asset-001" to:
      | field       | value        |
      | facility_id | facility-001 |
      | location_id | roof-level-3 |
    Then I receive a 403 Forbidden response

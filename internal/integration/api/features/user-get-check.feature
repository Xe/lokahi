Feature: A user can get a Check

  Background:
    Given a base stack
    Given an example check

  Scenario: normal get check
    As a User
    When I try to fetch the check
    Then the resulting check should have an ID
    Then tear everything down

  Scenario: cant get random check id
    As a User
    Given a random check ID
    When I try to fetch the check
    Then there was an error
    Then the check cannot be fetched
    Then tear everything down


Feature: challenge
  In order to check if challenges are correctly validated
  As a bot that uses the crypto api
  I need to be able to create and validate a challenge

  Scenario: challenge was successfully created
    Given a clean database
    When I send a request to create a challenge
    And I wait for the request to be processed
    Then the challenge should be created and valid

  Scenario: challenge was successfully created and validated
    Given a clean database
    Given a challenge that was previously created
    When I send a request to validate a challenge
    Then the challenge should be validated successfully
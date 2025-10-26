*** Settings ***
Library    SeleniumLibrary
Suite Setup    Open Browser To SauceDemo
Suite Teardown    Close Browser

*** Variables ***
${URL}              https://www.saucedemo.com/
${BROWSER}          edge
${USERNAME}         standard_user
${PASSWORD}         secret_sauce
${FIRSTNAME}        Pong
${LASTNAME}         Pp
${POSTALCODE}       12345

*** Keywords ***
Open Browser To SauceDemo
    Open Browser    ${URL}    ${BROWSER}
    Maximize Browser Window
    Set Selenium Speed    0.5s

Login To SauceDemo
    Input Text    id=user-name    ${USERNAME}
    Input Text    id=password     ${PASSWORD}
    Click Button  id=login-button
    Wait Until Page Contains Element    class=inventory_list

Add Product To Cart
    Click Button    id=add-to-cart-sauce-labs-backpack
    Click Element   class=shopping_cart_link
    Click Button    id=continue-shopping
    Click Button    id=add-to-cart-sauce-labs-bike-light
    Click Element   class=shopping_cart_link
    Wait Until Page Contains Element    class=cart_item
Checkout Order
    Click Button    id=checkout
    Input Text      id=first-name    ${FIRSTNAME}
    Input Text      id=last-name     ${LASTNAME}
    Input Text      id=postal-code   ${POSTALCODE}
    Click Button    id=continue
    # Wait Until Page Contains    Total: $43.18
    Click Button    id=finish
    Wait Until Page Contains    Thank you for your order!

*** Test Cases ***
E2E Purchase Flow
    Login To SauceDemo
    Add Product To Cart
    Checkout Order

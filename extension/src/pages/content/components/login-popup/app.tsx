import { useEffect, useState } from "react";
import { NextUIProvider, Button, Card, Text, Row, Tooltip } from '@nextui-org/react';
import log from "@src/utils/log"
import { ENDPOINTS, EXTENSION_NAME } from "@src/constants";

function optOut(email: string) {
  let toSet = {}
  toSet[email] = {optOut: true}
  chrome.storage.local.set(toSet);
}

export default function App() {
  const email = window['GLOBALS'][10]

  useEffect(() => {
    log("Login popup loaded");
  }, []);

  return (
    <NextUIProvider disableBaseline={true}>
      <Card variant="bordered" css={{ mw: "12cm" }}>
        <Card.Header>
          <Text b>Connect Account</Text>
        </Card.Header>
        <Card.Divider />
        <Card.Body>
          <Text>
            {email} is not connected to {EXTENSION_NAME}
          </Text>
        </Card.Body>
        <Card.Divider />
        <Card.Footer>
          <Row>
            <Button size="sm" onPress={() => {
              optOut(email);
              document.getElementById('lah-popup').remove();
            }} light>
              No thanks
            </Button>
            <Button size="sm" color="primary" onPress={() => window.open(ENDPOINTS.AUTH(email))}>
              Connect
            </Button>
          </Row>
        </Card.Footer>
      </Card>
    </NextUIProvider>
  )
}

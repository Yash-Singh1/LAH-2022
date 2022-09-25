import { useEffect, useState } from "react";
import {
  NextUIProvider,
  Button,
  Card,
  Text,
  Row,
  Tooltip,
  Avatar,
  Grid,
  Loading,
} from "@nextui-org/react";
import log from "@src/utils/log";
import { ENDPOINTS, EXTENSION_NAME, POPUP_KEYS } from "@src/constants";

interface Email {
  from_domain: string;
  sentiment: boolean;
  isSpa: boolean;
  organization: string;
}

const hide = () => document.getElementById(POPUP_KEYS.INFO).remove();

export default function App() {
  const [emailInfo, setEmailInfo] = useState(undefined);
  window["setEmailInfo"] = setEmailInfo;

  useEffect(() => {
    log("Info popup loaded");
  }, []);

  return (
    <NextUIProvider disableBaseline={true}>
      <Card variant="bordered" css={{ mw: "12cm" }}>
        {emailInfo == undefined ? (
          <Loading></Loading>
        ) : (
          <>
            <Card.Header
              css={{ boxSizing: "border-box", justifyContent: "space-between" }}
            >
              <div
                style={{
                  display: "flex",
                  justifyContent: "center",
                  alignItems: "center",
                }}
              >
                <Avatar
                  src={`https://www.google.com/s2/favicons?domain=${emailInfo[
                    "from_domain"
                  ].replace(/^google\.com$/, "www.google.com")}`}
                  squared
                  size="sm"
                />
                <Text b h3 css={{ marginLeft: "0.1rem" }}>
                  {emailInfo["from_domain"]}
                </Text>
              </div>
              <Text b>&times;</Text>
            </Card.Header>
            <Card.Divider />
            <Card.Body>
              <Text b h2>
                Analysis
              </Text>
              <Text b h3 color={!emailInfo["sentiment"] ? "success" : "error"}>
                Is this spam?: {!emailInfo["sentiment"] ? "👍" : "❌"}
              </Text>
              <Text b h3 color={!emailInfo.is_spam ? "success" : "error"}>
                Mood: {!emailInfo.is_spam ? "😊" : "😞"}
              </Text>
            </Card.Body>
            <Card.Divider />
            <Card.Footer>
              <Row>
                <Button color="primary" size="sm" onPress={hide}>
                  Okay! Thanks
                </Button>
              </Row>
            </Card.Footer>
          </>
        )}
      </Card>
    </NextUIProvider>
  );
}

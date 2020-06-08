import React from "react";
import { Box, Container, Typography } from "@material-ui/core";
import MuiLink from "@material-ui/core/Link";
import EmailIcon from "@material-ui/icons/Email";
import Copyright from "../src/Copyright";
import Deals from "../src/Deals";
import Spacer from "../src/Spacer";
import Header from "../src/Header";

export default function Index() {
  return (
    <Container>
      <Box my={4}>
        <Header />
        <Spacer />
        <Typography variant="subtitle2" gutterBottom>
          <EmailIcon color="primary" style={{ verticalAlign: "middle" }} />
          <MuiLink href="/subscribe">
            {` Want to get these deals in your inbox every morning? Click to
            subscribe to the newsletter! `}
          </MuiLink>
        </Typography>
        <Spacer />
        <Deals />
        <Spacer />
        <Copyright />
      </Box>
    </Container>
  );
}

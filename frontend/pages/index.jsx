import React from "react";
import EmailIcon from "@material-ui/icons/Email";
import Copyright from "../src/Copyright";
import Deals from "../src/Deals";
import Spacer from "../src/Spacer";
import Header from "../src/Header";
import Link from "next/link";
import Container from "@material-ui/core/Container";
import Box from "@material-ui/core/Box";
import Typography from "@material-ui/core/Typography";
import MuiLink from "@material-ui/core/Link";

export default function Index() {
  return (
    <Container>
      <Box my={4}>
        <Header />
        <Spacer />
        <Link href="/subscribe">
          <Typography variant="subtitle2" gutterBottom>
            <EmailIcon color="primary" style={{ verticalAlign: "middle" }} />
            <MuiLink>
              {` Want to get these deals in your inbox every morning? Click to
            subscribe to the newsletter! `}
            </MuiLink>
          </Typography>
        </Link>
        <Spacer />
        <Deals />
        <Spacer />
        <Copyright />
      </Box>
    </Container>
  );
}

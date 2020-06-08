import React from "react";
import { Container, Typography, Box } from "@material-ui/core";
import Copyright from "../src/Copyright";
import Deals from "../src/Deals";
import Spacer from "../src/Spacer";

export default function Index() {
  return (
    <Container>
      <Box my={4}>
        <Typography variant="h1" component="h1" gutterBottom>
          Hot Flag Deals
        </Typography>
        <Deals />
        <Spacer />
        <Copyright />
      </Box>
    </Container>
  );
}

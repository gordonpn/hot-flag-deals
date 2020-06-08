import React from "react";
import { Container, Typography, Box } from "@material-ui/core";
import Copyright from "../src/Copyright";
import Deals from "../src/Deals";
import Spacer from "../src/Spacer";
import { makeStyles } from "@material-ui/core/styles";

const useStyles = makeStyles({
  heading: {
    fontWeight: 400,
    background: "linear-gradient(180deg, rgba(255,255,255,0) 75%, #ff839c 75%)",
    display: "inline",
  },
});

export default function Index() {
  const classes = useStyles();
  return (
    <Container>
      <Box my={4}>
        <Typography
          variant="h1"
          component="h1"
          gutterBottom
          className={classes.heading}
        >
          Hot Flag Deals
        </Typography>
        <Spacer />
        <Deals />
        <Spacer />
        <Copyright />
      </Box>
    </Container>
  );
}

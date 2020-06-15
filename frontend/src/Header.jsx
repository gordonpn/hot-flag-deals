import React from "react";
import { makeStyles } from "@material-ui/core/styles";
import Typography from "@material-ui/core/Typography";

const useStyles = makeStyles({
  heading: {
    fontFamily: "'Asap Condensed', sans-serif",
    fontWeight: 500,
    background: "linear-gradient(180deg, rgba(255,255,255,0) 75%, #ff839c 75%)",
    display: "inline",
  },
});

export default function Header() {
  const classes = useStyles();
  return (
    <Typography
      variant="h1"
      component="h1"
      gutterBottom
      className={classes.heading}
    >
      Hot Flag Deals
    </Typography>
  );
}

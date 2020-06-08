import React from "react";
import { Typography } from "@material-ui/core";
import MuiLink from "@material-ui/core/Link";

export default function Copyright() {
  return (
    <Typography variant="body2" color="textSecondary" align="center">
      {"Copyright "}
      &copy;
      {` ${new Date().getFullYear()}`}
      {", Gordon Pham-Nguyen "}
      <br />
      <MuiLink
        color="inherit"
        href="http://gordon-pn.com"
        target="_blank"
        rel="noopener noreferrer"
      >
        gordon-pn.com
      </MuiLink>
      <br />
      <MuiLink
        color="inherit"
        href="https://github.com/gordonpn/hot-flag-deals"
        target="_blank"
        rel="noopener noreferrer"
      >
        Source Code on GitHub
      </MuiLink>
    </Typography>
  );
}

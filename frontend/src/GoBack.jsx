import Link from "next/link";
import Typography from "@material-ui/core/Typography";
import ArrowBackIcon from "@material-ui/icons/ArrowBack";
import MuiLink from "@material-ui/core/Link";
import React from "react";

export default function GoBack() {
  return (
    <Link href="/">
      <Typography color="primary" variant="subtitle2" gutterBottom>
        <ArrowBackIcon color="primary" style={{ verticalAlign: "middle" }} />
        <MuiLink>{` Go back to the deals. `}</MuiLink>
      </Typography>
    </Link>
  );
}

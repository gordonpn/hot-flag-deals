import { useRouter } from "next/router";
import React, { useEffect, useState } from "react";
import Container from "@material-ui/core/Container";
import Box from "@material-ui/core/Box";
import Header from "../src/Header";
import Spacer from "../src/Spacer";
import Copyright from "../src/Copyright";
import Typography from "@material-ui/core/Typography";
import GoBack from "../src/GoBack";
import * as Yup from "yup";
import axios from "axios";

export const schema = Yup.object().shape({
  email: Yup.string().email().lowercase().trim(),
});

export default function Confirm() {
  const router = useRouter();
  const [message, setMessage] = useState("");

  useEffect(() => {
    const { email } = router.query;
    if (email === undefined) {
      return;
    }
    schema
      .isValid({
        email,
      })
      .then((value) => {
        if (value) {
          axios
            .put("/api/v1/emails", { email })
            .then(() => {
              setMessage("Thank you for confirming your subscription!");
            })
            .catch(() => {
              setMessage("Something went wrong.");
            });
        } else {
          setMessage("Something went wrong.");
        }
      });
  }, [router.query]);

  return (
    <Container>
      <Box my={4}>
        <Header />
        <Spacer />
        <GoBack />
        <Spacer />
        <Typography align="center" variant="h2">
          {message}
        </Typography>
        <Spacer />
        <Copyright />
      </Box>
    </Container>
  );
}

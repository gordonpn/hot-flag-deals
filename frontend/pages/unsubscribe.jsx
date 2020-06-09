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

const schema = Yup.object().shape({
  email: Yup.string().email(),
});

export default function Unsubscribe() {
  const router = useRouter();
  const [message, setMessage] = useState("");

  useEffect(() => {
    const { email } = router.query;
    if (email !== undefined) {
      schema
        .isValid({
          email,
        })
        .then((value) => {
          if (value) {
            setMessage("You've been unsubscribed.");
            //  todo make call to backend with email
          } else {
            setMessage("Something went wrong.");
          }
        });
    } else {
      setMessage("Something went wrong.");
    }
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

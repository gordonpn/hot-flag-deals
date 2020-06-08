import Copyright from "../src/Copyright";
import {
  Container,
  Box,
  Button,
  LinearProgress,
  Typography,
} from "@material-ui/core";
import React, { useState } from "react";
import Header from "../src/Header";
import Spacer from "../src/Spacer";
import { Formik, Form, Field } from "formik";
import { TextField } from "formik-material-ui";
import { makeStyles } from "@material-ui/core/styles";
import MuiLink from "@material-ui/core/Link";
import ArrowBackIcon from "@material-ui/icons/ArrowBack";

const useStyles = makeStyles(() => ({
  alignItemsAndJustifyContent: {
    width: 720,
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
  },
  fields: {
    width: 360,
  },
}));

export default function Subscribe() {
  const classes = useStyles();
  const [submitted, setSubmitted] = useState(false);
  const [name, setName] = useState("");
  return (
    <Container>
      <Box my={4}>
        <Header />
        <Spacer />
        <Typography variant="subtitle2" gutterBottom>
          <ArrowBackIcon color="primary" style={{ verticalAlign: "middle" }} />
          <MuiLink href="/">{` Go back to the deals. `}</MuiLink>
        </Typography>
        {submitted && (
          <>
            <Spacer />
            <Typography color="primary" variant="subtitle1" gutterBottom>
              {`Thanks for subscribing ${name}!`}
            </Typography>
            <Typography color="primary" variant="subtitle1" gutterBottom>
              Please check your inbox to confirm your subscription.
            </Typography>
          </>
        )}
        <Spacer />
        <div className={classes.alignItemsAndJustifyContent}>
          <Formik
            initialValues={{
              name: "",
              email: "",
            }}
            validate={(values) => {
              const errors = {};
              if (!values.email) {
                errors.email = "Required";
              } else if (
                !/^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}$/i.test(values.email)
              ) {
                errors.email = "Invalid email address";
              } else if (
                !/^[a-zA-Z]+(([',. -][a-zA-Z ])?[a-zA-Z]*)*$/i.test(values.name)
              ) {
                errors.name = "Invalid name";
              }
              return errors;
            }}
            onSubmit={(values, { setSubmitting }) => {
              setSubmitted(true);
              setName(values.name);
              setSubmitting(false);
            }}
          >
            {({ submitForm, isSubmitting }) => (
              <Form>
                <Field
                  component={TextField}
                  type="text"
                  label="Name"
                  name="name"
                  className={classes.fields}
                />
                <br />
                <Field
                  component={TextField}
                  name="email"
                  type="email"
                  label="Email"
                  className={classes.fields}
                />
                {isSubmitting && <LinearProgress />}
                <br />
                <Spacer />
                <Button
                  variant="contained"
                  color="primary"
                  disabled={isSubmitting}
                  onClick={submitForm}
                >
                  Subscribe
                </Button>
              </Form>
            )}
          </Formik>
        </div>
        <Spacer />
        <Copyright />
      </Box>
    </Container>
  );
}

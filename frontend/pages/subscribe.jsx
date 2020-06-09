import Copyright from "../src/Copyright";
import {
  Box,
  Button,
  Container,
  LinearProgress,
  Typography,
} from "@material-ui/core";
import React, { useState } from "react";
import Header from "../src/Header";
import Spacer from "../src/Spacer";
import { Field, Form, Formik } from "formik";
import { TextField } from "formik-material-ui";
import { makeStyles } from "@material-ui/core/styles";
import MuiLink from "@material-ui/core/Link";
import ArrowBackIcon from "@material-ui/icons/ArrowBack";
import ReCAPTCHA from "react-google-recaptcha";
import * as Yup from "yup";

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
              recaptcha: "",
            }}
            validationSchema={Yup.object().shape({
              name: Yup.string().matches(
                /^[a-zA-Z]+(([',. -][a-zA-Z ])?[a-zA-Z]*)*$/,
                "Invalid name"
              ),
              email: Yup.string()
                .email("Invalid email")
                .required("Email required"),
              recaptcha: Yup.string().required("recaptcha required"),
            })}
            onSubmit={(values, { setSubmitting }) => {
              setSubmitted(true);
              setName(values.name);
              setSubmitting(false);
            }}
          >
            {({ errors, touched, submitForm, isSubmitting, setFieldValue }) => (
              <Form>
                <Field
                  component={TextField}
                  type="text"
                  label="Name"
                  name="name"
                  variant="filled"
                  className={classes.fields}
                />
                <br />
                <Field
                  component={TextField}
                  name="email"
                  type="email"
                  label="Email"
                  variant="filled"
                  className={classes.fields}
                />
                {isSubmitting && <LinearProgress />}
                <br />
                <Spacer />
                <Field name="recaptcha" style={{ display: "none" }} />
                <ReCAPTCHA
                  sitekey="6LdAlQEVAAAAAGKrXMMe55XXlcknuswppK9xXpUI"
                  onChange={(value) => {
                    setFieldValue("recaptcha", value);
                  }}
                />
                {errors.recaptcha && touched.recaptcha ? (
                  <Typography color="secondary">{errors.recaptcha}</Typography>
                ) : null}
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

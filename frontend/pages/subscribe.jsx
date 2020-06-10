import Copyright from "../src/Copyright";
import React, { useState } from "react";
import Header from "../src/Header";
import Spacer from "../src/Spacer";
import { Field, Form, Formik } from "formik";
import { TextField } from "formik-material-ui";
import { makeStyles } from "@material-ui/core/styles";
import ReCAPTCHA from "react-google-recaptcha";
import * as Yup from "yup";
import Container from "@material-ui/core/Container";
import Box from "@material-ui/core/Box";
import Typography from "@material-ui/core/Typography";
import LinearProgress from "@material-ui/core/LinearProgress";
import Button from "@material-ui/core/Button";
import GoBack from "../src/GoBack";

const useStyles = makeStyles(() => ({
  alignItemsAndJustifyContent: {
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
  },
  fields: {
    width: 360,
  },
}));

const signUpSchema = Yup.object().shape({
  name: Yup.string()
    .matches(/^[a-zA-Z]+(([',. -][a-zA-Z ])?[a-zA-Z]*)*$/, "Invalid name")
    .trim(),
  email: Yup.string()
    .email("Invalid email")
    .required("Email required")
    .lowercase()
    .trim(),
  recaptcha: Yup.string().required("recaptcha required").ensure(),
});

export default function Subscribe() {
  const classes = useStyles();
  const [submitted, setSubmitted] = useState(false);
  const [name, setName] = useState("");

  return (
    <Container>
      <Box my={4}>
        <Header />
        <Spacer />
        <GoBack />
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
            validationSchema={signUpSchema}
            onSubmit={(values, { setSubmitting }) => {
              setSubmitted(true);
              setName(values.name);
              setSubmitting(false);
              //  todo make call to backend with email
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

import {
  Box,
  AppBar,
  Toolbar,
  IconButton,
  Typography,
  Tooltip,
  Container,
  TextField,
  Button,
  Snackbar,
} from "@mui/material";
import React, { useCallback, useContext, useEffect, useState } from "react";
import { AuthContext } from "../contexts/auth";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import { useNavigate } from "react-router-dom";
import { HOME_ROUTE } from "../config/routes";
import UserAvatar from "../components/UserAvatar";
import { useFormik } from "formik";
import theme from "../config/theme";
import LoadingSwitcher from "../components/LoadingSwitcher";
import CloseIcon from "@mui/icons-material/Close";
import { createShweet } from "../services/shweets";
import ApiError from "../models/errors/ApiError";

interface Values {
  message: string;
}

const NewShweetPage = () => {
  const navigate = useNavigate();
  const { user } = useContext(AuthContext);
  const [snackbarOpen, setSnackbarOpen] = useState(false);

  const handleClose = (
    event: React.SyntheticEvent | Event,
    reason?: string
  ) => {
    if (reason === "clickaway") {
      return;
    }

    setSnackbarOpen(false);
  };

  const formik = useFormik<Values>({
    initialValues: { message: "" },
    onSubmit: async (values, { setSubmitting, setErrors, setStatus }) => {
      setStatus(undefined);

      const result = await createShweet(values);
      if (result instanceof ApiError) {
        setErrors(result.getError());
        setStatus(result.getFormattedStatus());
      } else {
        // TODO: on success, animate a screen of success and the go to the refreshed home page
        navigate(-1);
      }

      setSubmitting(false);
    },
  });

  const keydownHandler = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === "Enter" && e.ctrlKey) {
        formik.submitForm();
      }
    },
    [formik]
  );

  useEffect(() => {
    window.addEventListener("keydown", keydownHandler);
    return () => window.removeEventListener("keydown", keydownHandler);
  }, [keydownHandler]);

  if (user === undefined) {
    return <></>;
  }

  return (
    <>
      <Box>
        <AppBar position="static">
          <Toolbar>
            <Tooltip title="Go back">
              <IconButton
                edge="start"
                color="inherit"
                aria-label="menu"
                sx={{ mr: 2 }}
                onClick={() => navigate(HOME_ROUTE)}
              >
                <ArrowBackIcon />
              </IconButton>
            </Tooltip>
          </Toolbar>
        </AppBar>
      </Box>

      <Box
        sx={{
          margin: 4,
          display: "flex",
          flexDirection: "column",
          alignItems: "start",
        }}
      >
        <Container component="main" maxWidth="xs">
          <Box mb={5}>
            <Typography variant="h6" gutterBottom>
              Shweet something!
            </Typography>
            <Typography variant="body1">
              You can shweet anything, the shittier the better! No one is going
              to read it anyway.
            </Typography>
          </Box>

          <Box display="flex" alignItems="center">
            <Box mr={2}>
              <UserAvatar user={user} />
            </Box>
            <Typography>{user.name}</Typography>
            <Button
              sx={{ marginLeft: "auto" }}
              onClick={() => setSnackbarOpen(true)}
            >
              Switch
            </Button>
          </Box>

          <Box
            component="form"
            onSubmit={formik.handleSubmit}
            noValidate
            sx={{ mt: 2 }}
          >
            <TextField
              id="message"
              name="message"
              label="What are you thinking about?"
              fullWidth
              required
              autoFocus
              multiline
              minRows={4}
              value={formik.values.message}
              onChange={formik.handleChange}
              error={formik.touched.message && Boolean(formik.errors.message)}
              helperText={formik.touched.message && formik.errors.message}
            />

            <Box mt={1} minHeight={24}>
              <Typography
                style={{ color: theme.palette.error.main }}
                align="center"
              >
                {formik.status}
              </Typography>
            </Box>

            <Button
              fullWidth
              variant="contained"
              sx={{ mt: 1, mb: 2 }}
              type="submit"
              disabled={formik.isSubmitting}
            >
              <LoadingSwitcher
                color="secondary"
                size={24}
                loading={formik.isSubmitting}
              >
                Shweet it
              </LoadingSwitcher>
            </Button>
          </Box>
        </Container>
        <Snackbar
          open={snackbarOpen}
          autoHideDuration={6000}
          onClose={handleClose}
          message="Not implemented yet! Shweet a complaint to @shwitter"
          action={
            <IconButton
              size="small"
              aria-label="close"
              color="inherit"
              onClick={handleClose}
            >
              <CloseIcon fontSize="small" />
            </IconButton>
          }
        />
      </Box>
    </>
  );
};

export default NewShweetPage;

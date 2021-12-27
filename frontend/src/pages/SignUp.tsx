import * as React from "react";
import Avatar from "@mui/material/Avatar";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Link from "@mui/material/Link";
import Grid from "@mui/material/Grid";
import Box from "@mui/material/Box";
import LockOutlinedIcon from "@mui/icons-material/LockOutlined";
import Typography from "@mui/material/Typography";
import Container from "@mui/material/Container";
import { useFormik } from "formik";
import { Link as RouterLink, useNavigate } from "react-router-dom";
import { AuthContext, AuthStatus } from "../contexts/auth";
import { createUser } from "../services/auth";
import ApiError from "../models/errors/ApiError";
import theme from "../config/theme";
import LoadingSwitcher from "../components/LoadingSwitcher";
import { HOME_ROUTE } from "../config/routes";

interface Values {
  name: string;
  username: string;
  email: string;
  password: string;
}

function Copyright(props: any) {
  return (
    <Typography
      variant="body2"
      color="text.secondary"
      align="center"
      {...props}
    >
      {"Copyright © "}
      <Link color="inherit" href="https://mui.com/">
        Your Website
      </Link>{" "}
      {new Date().getFullYear()}
      {"."}
    </Typography>
  );
}

export default function SignUpPage() {
  const navigate = useNavigate();
  const { authStatus, authLogin } = React.useContext(AuthContext);

  React.useEffect(() => {
    if (authStatus === AuthStatus.Authenticated) {
      navigate(HOME_ROUTE);
    }
  }, [authStatus]); // eslint-disable-line react-hooks/exhaustive-deps

  const formik = useFormik<Values>({
    initialValues: { username: "", name: "", email: "", password: "" },
    onSubmit: async (values, { setSubmitting, setErrors, setStatus }) => {
      setStatus(undefined);

      const result = await createUser(values);
      if (result instanceof ApiError) {
        setErrors(result.getError());
        setStatus(result.getFormattedStatus());
      } else {
        // success!
        await authLogin(values.username, values.password);
      }

      setSubmitting(false);
    },
  });

  return (
    <Container component="main" maxWidth="xs">
      <Box
        sx={{
          marginTop: 8,
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
        }}
      >
        <Avatar sx={{ m: 1, bgcolor: "secondary.main" }}>
          <LockOutlinedIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          Sign up
        </Typography>
        <Box component="form" onSubmit={formik.handleSubmit} sx={{ mt: 3 }}>
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <TextField
                autoComplete="username"
                name="username"
                required
                fullWidth
                id="username"
                label="Username"
                autoFocus
                value={formik.values.username}
                onChange={formik.handleChange}
                error={
                  formik.touched.username && Boolean(formik.errors.username)
                }
                helperText={formik.touched.username && formik.errors.username}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                required
                fullWidth
                id="name"
                label="Name"
                name="name"
                autoComplete="name"
                value={formik.values.name}
                onChange={formik.handleChange}
                error={formik.touched.name && Boolean(formik.errors.name)}
                helperText={formik.touched.name && formik.errors.name}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                required
                fullWidth
                id="email"
                label="Email Address"
                name="email"
                autoComplete="email"
                value={formik.values.email}
                onChange={formik.handleChange}
                error={formik.touched.email && Boolean(formik.errors.email)}
                helperText={formik.touched.email && formik.errors.email}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                required
                fullWidth
                name="password"
                label="Password"
                type="password"
                id="password"
                autoComplete="new-password"
                value={formik.values.password}
                onChange={formik.handleChange}
                error={
                  formik.touched.password && Boolean(formik.errors.password)
                }
                helperText={formik.touched.password && formik.errors.password}
              />
            </Grid>
          </Grid>

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
              Sign Up
            </LoadingSwitcher>
          </Button>
          <Grid container justifyContent="flex-end">
            <Grid item>
              <Link href="#" variant="body2"></Link>
              <RouterLink to="/sign-in">
                <Link variant="body2"> Already have an account? Sign in</Link>
              </RouterLink>
            </Grid>
          </Grid>
        </Box>
      </Box>
      <Copyright sx={{ mt: 5 }} />
    </Container>
  );
}

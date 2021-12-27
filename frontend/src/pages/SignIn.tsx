import * as React from "react";
import Avatar from "@mui/material/Avatar";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Grid from "@mui/material/Grid";
import Box from "@mui/material/Box";
import LockOutlinedIcon from "@mui/icons-material/LockOutlined";
import Typography from "@mui/material/Typography";
import Container from "@mui/material/Container";
import { AuthContext, AuthStatus } from "../contexts/auth";
import { useFormik } from "formik";
import theme from "../config/theme";
import { Link as RouterLink, useNavigate } from "react-router-dom";
import LoadingSwitcher from "../components/LoadingSwitcher";
import { Link } from "@mui/material";
import { HOME_ROUTE, SIGN_UP_ROUTE } from "../config/routes";

interface Values {
  username: string;
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

export default function SignInPage() {
  const navigate = useNavigate();
  const { authStatus, authLogin } = React.useContext(AuthContext);

  React.useEffect(() => {
    if (authStatus === AuthStatus.Authenticated) {
      navigate(HOME_ROUTE);
    }
  }, [authStatus]); // eslint-disable-line react-hooks/exhaustive-deps

  const formik = useFormik<Values>({
    initialValues: { username: "", password: "" },
    onSubmit: async (values, { setSubmitting, setErrors, setStatus }) => {
      setStatus(undefined);

      const result = await authLogin(values.username, values.password);
      if (result) {
        setErrors(result.getError());
        setStatus(result.getFormattedStatus());
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
          Sign in
        </Typography>
        <Box
          component="form"
          onSubmit={formik.handleSubmit}
          noValidate
          sx={{ mt: 1 }}
        >
          <TextField
            margin="normal"
            id="username"
            name="username"
            label="Username"
            fullWidth
            required
            autoFocus
            value={formik.values.username}
            onChange={formik.handleChange}
            error={formik.touched.username && Boolean(formik.errors.username)}
            helperText={formik.touched.username && formik.errors.username}
          />

          <TextField
            margin="normal"
            id="password"
            name="password"
            label="Password"
            fullWidth
            required
            type="password"
            autoComplete="current-password"
            value={formik.values.password}
            onChange={formik.handleChange}
            error={formik.touched.password && Boolean(formik.errors.password)}
            helperText={formik.touched.password && formik.errors.password}
          />
          {/* <FormControlLabel
              control={<Checkbox value="remember" color="primary" />}
              label="Remember me"
            /> */}
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
              Sign In
            </LoadingSwitcher>
          </Button>
          <Grid container>
            <Grid item xs>
              <RouterLink to="/reset-password">
                <Link variant="body2">Forgot password?</Link>
              </RouterLink>
            </Grid>
            <Grid item>
              <RouterLink to={SIGN_UP_ROUTE}>
                <Link variant="body2">Don't have an account? Sign Up</Link>
              </RouterLink>
            </Grid>
          </Grid>
        </Box>
      </Box>
      <Copyright sx={{ mt: 8, mb: 4 }} />
    </Container>
  );
}

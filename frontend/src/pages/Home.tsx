import { Container, Box, Typography, Button } from "@mui/material";
import { useContext, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { AuthContext, AuthStatus } from "../contexts/auth";

const HomePage = () => {
  const navigate = useNavigate();
  const { authStatus, authLogout } = useContext(AuthContext);

  useEffect(() => {
    if (authStatus === AuthStatus.NotAuthenticated) {
      navigate("/sign-in");
    }
  }, [authStatus]);

  return (
    <Container maxWidth="sm">
      <Box py={4}>
        <main>
          <Typography variant="h5" gutterBottom align="center">
            Welcome to Shwitter!
          </Typography>
          <Button onClick={authLogout}>Logout</Button>
        </main>
      </Box>
    </Container>
  );
};

export default HomePage;

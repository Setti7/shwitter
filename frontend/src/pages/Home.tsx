import { Container, Box } from "@mui/material";
import { useContext, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import Header from "../components/Header";
import { AuthContext, AuthStatus } from "../contexts/auth";

const HomePage = () => {
  const navigate = useNavigate();
  const { user, authStatus } = useContext(AuthContext);

  useEffect(() => {
    if (authStatus === AuthStatus.NotAuthenticated) {
      navigate("/sign-in");
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [authStatus]);

  if (user === undefined) {
    return <></>;
  }

  return (
    <>
      <Header />
      <Container component="main">
        <Box
          sx={{
            marginTop: 4,
            display: "flex",
            flexDirection: "column",
            alignItems: "stretch",
          }}
        >
          {/* TODO: 
          [ ] Add the floating tweet button 
          [ ] Add timeline
          [ ] Add userline
          */}

        </Box>
      </Container>
    </>
  );
};

export default HomePage;

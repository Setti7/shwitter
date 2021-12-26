import { Box, Container, Typography } from "@mui/material";
import React from "react";

const Header: React.FC = () => {
  return (
    <Container maxWidth="sm">
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Header
        </Typography>
      </Box>
    </Container>
  );
};

export default Header;

import { Box, CircularProgress } from "@mui/material";
import React from "react";

interface Props {
  loading: boolean;
  size?: number;
  color?: "primary" | "secondary" | "inherit";
}
const LoadingSwitcher: React.FC<Props> = ({
  loading,
  children,
  size,
  color,
}) => {
  return (
    <Box
      style={{ position: "relative" }}
      display="flex"
      justifyContent="center"
      alignItems="center"
    >
      <Box
        style={{ position: "absolute" }}
        visibility={!loading ? "hidden" : "visible"}
        width={size}
        height={size}
      >
        <CircularProgress size={size} color={color} />
      </Box>
      <Box visibility={loading ? "hidden" : "visible"}>{children}</Box>
    </Box>
  );
};

export default LoadingSwitcher;

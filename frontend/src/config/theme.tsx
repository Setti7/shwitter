import { createTheme } from "@mui/material";

const theme = createTheme({
  palette: {
    mode: "dark",
    // primary: {
    //   light: "#757ce8",
    //   main: "#3f50b5",
    //   dark: "#002884",
    //   contrastText: "#fff",
    // },
    // secondary: {
    //   light: "#ff7961",
    //   main: "#f44336",
    //   dark: "#ba000d",
    //   contrastText: "#000",
    // },
  },
  typography: {
    fontFamily: '"Signika", "Helvetica", "Arial", sans-serif',
    caption: {
      color: "#6E767D",
      fontSize: "0.90rem"
    },
  },
});

export const fabBlackStyle = {
  color: "common.white",
  bgcolor: "#000",
  borderRadius: "100%",
  border: "1px solid #666",
  "&:hover": {
    bgcolor: "#333",
  },
};

export default theme;

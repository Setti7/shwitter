import * as React from "react";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import IconButton from "@mui/material/IconButton";
import { AuthContext } from "../contexts/auth";
import { Tooltip } from "@mui/material";
import LogoutIcon from "@mui/icons-material/Logout";
import UserAvatar from "./UserAvatar";

const Header = () => {
  const { user, authLogout } = React.useContext(AuthContext);

  if (user === undefined) {
    return <></>;
  }

  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static">
        <Toolbar>
          <IconButton
            edge="start"
            color="inherit"
            aria-label="menu"
            sx={{ mr: 2 }}
          >
            <UserAvatar user={user} />
          </IconButton>

          <Typography variant="h6" sx={{ flexGrow: 1 }}>
            Home
          </Typography>
          <Tooltip title="logout">
            <IconButton
              color="inherit"
              aria-label="logout"
              onClick={authLogout}
            >
              <LogoutIcon />
            </IconButton>
          </Tooltip>
        </Toolbar>
      </AppBar>
    </Box>
  );
};

export default Header;

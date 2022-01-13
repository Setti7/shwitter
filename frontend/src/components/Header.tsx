import * as React from "react";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import IconButton from "@mui/material/IconButton";
import { AuthContext, AuthStatus } from "../contexts/auth";
import { Tooltip } from "@mui/material";
import LogoutIcon from "@mui/icons-material/Logout";
import LockOpenIcon from "@mui/icons-material/LockOpen";
import UserAvatar from "./UserAvatar";
import { Link as RouterLink, useNavigate } from "react-router-dom";
import { HOME_ROUTE, SIGN_IN_ROUTE } from "../config/routes";

const Header = () => {
  const { user, authLogout, authStatus } = React.useContext(AuthContext);
  const navigate = useNavigate();

  if (user === undefined) {
    return (
      <Box sx={{ flexGrow: 1 }}>
        <AppBar position="static">
          <Toolbar>
            <Typography variant="h6" sx={{ flexGrow: 1 }}>
              <RouterLink
                to={HOME_ROUTE}
                style={{ textDecoration: "none", color: "white" }}
              >
                Home
              </RouterLink>
            </Typography>

            {authStatus === AuthStatus.Authenticated ? (
              <Tooltip title="Logout">
                <IconButton
                  color="inherit"
                  aria-label="Logout"
                  onClick={authLogout}
                >
                  <LogoutIcon />
                </IconButton>
              </Tooltip>
            ) : (
              <Tooltip title="Sign in">
                <IconButton
                  color="inherit"
                  aria-label="Sign in"
                  onClick={() => navigate(SIGN_IN_ROUTE)}
                >
                  <LockOpenIcon />
                </IconButton>
              </Tooltip>
            )}
          </Toolbar>
        </AppBar>
      </Box>
    );
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
            <RouterLink
              to={HOME_ROUTE}
              style={{ textDecoration: "none", color: "white" }}
            >
              Home
            </RouterLink>
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

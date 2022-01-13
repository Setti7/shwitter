import { Box, Typography } from "@mui/material";
import { FC } from "react";
import User from "../models/user";

const UserDetails: FC<{ user: User }> = ({ user }) => {
  return (
    <>
      <Box display="flex" flexDirection="column" justifyContent="start">
        <Typography>{user.name}</Typography>
        <Box display="flex" alignItems="baseline">
          <Typography sx={{ flexGrow: 1 }} variant="caption">
            @{user.username}
          </Typography>
           {/* TODO: humanize user joined date */}
          <Typography variant="caption">Joined 11 months ago</Typography>
        </Box>

        <Box mt={1} />
        <Typography variant="body2">{user.bio}</Typography>

        <Box mt={1} display="flex" alignItems="end">
          <Box display="flex" flexGrow={1}>
            <Typography>
              XXX <Typography variant="caption">following</Typography>
            </Typography>
            <Box ml={2} />
            <Typography>
              XXX <Typography variant="caption">followers</Typography>
            </Typography>
          </Box>
        </Box>
      </Box>
    </>
  );
};

export default UserDetails;

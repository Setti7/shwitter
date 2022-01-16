import { Box, Typography } from "@mui/material";
import { FC } from "react";
import { UserProfile } from "../models/user";

const UserDetails: FC<{ userProfile: UserProfile }> = ({ userProfile: user }) => {
  return (
    <>
      <Box display="flex" flexDirection="column" justifyContent="start" m={2}>
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
              {user.friends_count}{" "}
              <Typography variant="caption">following</Typography>
            </Typography>
            <Box ml={2} />
            <Typography>
              {user.followers_count}{" "}
              <Typography variant="caption">followers</Typography>
            </Typography>
          </Box>
        </Box>
      </Box>
    </>
  );
};

export default UserDetails;

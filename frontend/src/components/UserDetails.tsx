import { Box, Typography } from "@mui/material";
import dayjs from "dayjs";
import { FC } from "react";
import { UserProfile } from "../models/user";

const UserDetails: FC<{ userProfile: UserProfile }> = ({ userProfile }) => {
  const now = dayjs();
  const joinedAtHumanized = dayjs
    .duration(now.diff(dayjs(userProfile.joined_at)))
    .humanize();

  return (
    <>
      <Box display="flex" flexDirection="column" justifyContent="start" m={2}>
        <Typography>{userProfile.name}</Typography>
        <Box display="flex" alignItems="baseline">
          <Typography sx={{ flexGrow: 1 }} variant="caption">
            @{userProfile.username}
          </Typography>
          <Typography variant="caption">
            Joined {joinedAtHumanized} ago
          </Typography>
        </Box>

        <Box mt={1} />
        <Typography variant="body2">{userProfile.bio}</Typography>

        {/* TODO: add page where user can list its followers/friends */}
        <Box mt={1} display="flex" alignItems="end">
          <Box display="flex" flexGrow={1}>
            <Typography>
              {userProfile.friends_count}{" "}
              <Typography variant="caption">following</Typography>
            </Typography>
            <Box ml={2} />
            <Typography>
              {userProfile.followers_count}{" "}
              <Typography variant="caption">followers</Typography>
            </Typography>
          </Box>
        </Box>
      </Box>
    </>
  );
};

export default UserDetails;

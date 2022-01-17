import { Box, Typography } from "@mui/material";
import dayjs from "dayjs";
import { FC } from "react";
import { Link } from "react-router-dom";
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

        <Box mt={1} display="flex" alignItems="end">
          <Box display="flex" flexGrow={1}>
            <Link
              to={"/user/" + userProfile.id + "/friends"}
              style={{ color: "inherit", textDecoration: "none" }}
            >
              <Typography>
                {userProfile.friends_count}{" "}
                <Typography variant="caption">following</Typography>
              </Typography>
            </Link>

            <Box ml={2} />
            <Link
              to={"/user/" + userProfile.id + "/followers"}
              style={{ color: "inherit", textDecoration: "none" }}
            >
              <Typography>
                {userProfile.followers_count}{" "}
                <Typography variant="caption">followers</Typography>
              </Typography>
            </Link>
          </Box>
        </Box>
      </Box>
    </>
  );
};

export default UserDetails;

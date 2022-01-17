import { Box, Typography } from "@mui/material";
import { FC } from "react";
import UserAvatar from "./UserAvatar";
import dayjs from "dayjs";
import { Link as RouterLink } from "react-router-dom";
import { FriendOrFollower } from "../models/user";

interface Props {
  friendOrFollower: FriendOrFollower;
}

const FriendOrFollowerCard: FC<Props> = ({ friendOrFollower }) => {
  return (
    <Box display="flex" flexDirection="row">
      <Box mr={2}>
        <UserAvatar user={friendOrFollower} />
      </Box>

      <Box
        display="flex"
        flexDirection="column"
        alignItems="stretch"
        flexGrow={1}
      >
        <RouterLink
          to={"/user/" + friendOrFollower.id}
          style={{ textDecoration: "none", color: "inherit" }}
        >
          <Box display="flex">
            <Box display="flex" flexGrow={1} alignItems="baseline">
              <Typography>{friendOrFollower.name}</Typography>
              <Typography sx={{ marginLeft: 1 }} variant="caption">
                @{friendOrFollower.username}
              </Typography>
            </Box>
            <Typography variant="caption" justifySelf="end">
              since {dayjs(friendOrFollower.since).format("DD MMM YYYY")}
            </Typography>
          </Box>
          <Typography variant="body2">{friendOrFollower.bio}</Typography>
        </RouterLink>
      </Box>
    </Box>
  );
};

export default FriendOrFollowerCard;

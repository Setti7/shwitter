import { Box, Typography } from "@mui/material";
import { FC } from "react";
import Shweet from "../models/shweet";
import UserAvatar from "./UserAvatar";
import dayjs from "dayjs";
import { Link as RouterLink } from "react-router-dom";

interface Props {
  shweet: Shweet;
}

const ShweetCard: FC<Props> = ({ shweet }) => {
  const now = dayjs();
  const createdAtHumanized = dayjs
    .duration(now.diff(dayjs(shweet.created_at)))
    .humanize();

  return (
    <Box display="flex" flexDirection="row">
      <Box mr={2}>
        <UserAvatar user={shweet.user} />
      </Box>

      <Box
        display="flex"
        flexDirection="column"
        alignItems="stretch"
        flexGrow={1}
      >
        <RouterLink
          to={"/shweet/" + shweet.id}
          style={{ textDecoration: "none", color: "inherit" }}
        >
          <Box display="flex">
            <Box display="flex" flexGrow={1} alignItems="baseline">
              <Typography>{shweet.user.name}</Typography>
              <Typography sx={{ marginLeft: 1 }} variant="caption">
                @{shweet.user.username}
              </Typography>
            </Box>
            <Typography variant="caption" justifySelf="end">
              {createdAtHumanized}
            </Typography>
            {/* TODO: 
            [X] Add likes
            [ ] Add reshweet
            [X] Add share button
            [ ] Maybe add responses
            [ ] List users who liked/reshweeted
            [X] Add click to go to details
            [X] Add click to go to profile
            */}
          </Box>
          <Typography variant="body2">{shweet.message}</Typography>
        </RouterLink>
      </Box>
    </Box>
  );
};

export default ShweetCard;

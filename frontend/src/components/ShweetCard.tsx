import { Box, IconButton, Tooltip, Typography } from "@mui/material";
import { FC, useState } from "react";
import { ShweetDetails } from "../models/shweet";
import UserAvatar from "./UserAvatar";
import dayjs from "dayjs";
import { Link as RouterLink } from "react-router-dom";
import CachedIcon from "@mui/icons-material/Cached";
import FavoriteBorderIcon from "@mui/icons-material/FavoriteBorder";
import FavoriteIcon from "@mui/icons-material/Favorite";
import ChatBubbleOutlineIcon from "@mui/icons-material/ChatBubbleOutline";
import { likeShweet } from "../services/shweets";

interface Props {
  initialShweet: ShweetDetails;
}

const ShweetCard: FC<Props> = ({ initialShweet }) => {
  const now = dayjs();
  const createdAtHumanized = dayjs
    .duration(now.diff(dayjs(initialShweet.created_at)))
    .humanize();

  const [shweet, setShweet] = useState(initialShweet);

  const likeOrUnlike = async () => {
    const likeIncrement = shweet.liked ? -1 : 1;

    setShweet({
      ...shweet,
      liked: !shweet.liked,
      like_count: shweet.like_count + likeIncrement,
    });
    await likeShweet(shweet.id);
  };

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
          </Box>
          <Typography variant="body2">{shweet.message}</Typography>
        </RouterLink>

        <Box
          mt={1}
          mr={2}
          display="flex"
          flexGrow={1}
          justifyContent="space-between"
        >
          <Box>
            <Tooltip title="reshweet">
              <IconButton size="small">
                <CachedIcon fontSize="inherit" />
              </IconButton>
            </Tooltip>
            <Typography variant="caption">{shweet.reshweet_count}</Typography>
          </Box>

          <Box ml={2}>
            <Tooltip title="like">
              <IconButton size="small" onClick={likeOrUnlike}>
                {shweet.liked ? (
                  <FavoriteIcon fontSize="inherit" color="error" />
                ) : (
                  <FavoriteBorderIcon fontSize="inherit" />
                )}
              </IconButton>
            </Tooltip>
            <Typography variant="caption">{shweet.like_count}</Typography>
          </Box>

          <Box ml={2}>
            <Tooltip title="comment">
              <IconButton size="small">
                <ChatBubbleOutlineIcon fontSize="inherit" />
              </IconButton>
            </Tooltip>
            <Typography variant="caption">{shweet.comment_count}</Typography>
          </Box>
        </Box>
      </Box>
    </Box>
  );
};

export default ShweetCard;

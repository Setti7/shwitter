import { Box, Typography, Container } from "@mui/material";
import { FC } from "react";
import Shweet from "../models/shweet";
import UserAvatar from "./UserAvatar";

interface Props {
  shweet: Shweet;
}

const ShweetCard: FC<Props> = ({ shweet }) => {
  return (
    <Container maxWidth="xs">
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
          <Box display="flex" alignItems="baseline">
            <Typography>{shweet.user.name}</Typography>
            <Typography sx={{ marginLeft: 1 }} variant="caption">
              @{shweet.user.username}
            </Typography>
            {/* TODO: 
            [ ] Add post time 
            [ ] Add likes
            [ ] Add reshweet
            [ ] Add share button
            [ ] Maybe add responses
            [ ] Add click to go to details
            [ ] Add click to go to profile
            */}
          </Box>

          <Typography variant="body2">{shweet.message}</Typography>
        </Box>
      </Box>
    </Container>
  );
};

export default ShweetCard;

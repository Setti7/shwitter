import {
  Box,
  AppBar,
  Toolbar,
  IconButton,
  Typography,
  Container,
  Divider,
  Tooltip,
  Snackbar,
} from "@mui/material";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import { ShweetDetails } from "../models/shweet";
import { useCallback, useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { HOME_ROUTE } from "../config/routes";
import { getShweetDetails, likeShweet } from "../services/shweets";
import ApiError from "../models/errors/ApiError";
import dayjs from "dayjs";
import UserAvatar from "../components/UserAvatar";
import ShareIcon from "@mui/icons-material/Share";
import CachedIcon from "@mui/icons-material/Cached";
import FavoriteBorderIcon from "@mui/icons-material/FavoriteBorder";
import FavoriteIcon from "@mui/icons-material/Favorite";
import ChatBubbleOutlineIcon from "@mui/icons-material/ChatBubbleOutline";
import CloseIcon from "@mui/icons-material/Close";

const ShweetDetailsPage = () => {
  const navigate = useNavigate();
  const [shweet, setShweet] = useState<ShweetDetails | undefined>();
  const [error, setError] = useState<String | undefined>();
  const { shweetId } = useParams();
  const [snackbarOpen, setSnackbarOpen] = useState(false);

  useEffect(() => {
    const getData = async () => {
      if (!shweetId) {
        navigate(HOME_ROUTE);
        return;
      }

      const result = await getShweetDetails(shweetId);

      if (result instanceof ApiError) {
        setError(result.getFormattedStatus());
      } else {
        setShweet(result);
      }
    };

    getData();
  }, [shweetId]);

  const likeOrUnlike = async () => {
    if (shweetId === undefined || shweet === undefined) {
      return;
    }

    let likeIncrement = 1;
    if (shweet.liked) {
      likeIncrement = -1;
    }

    setShweet({
      ...shweet,
      liked: !shweet.liked,
      like_count: shweet.like_count + likeIncrement,
    });
    await likeShweet(shweetId);
  };

  const handleClose = (
    event: React.SyntheticEvent | Event,
    reason?: string
  ) => {
    if (reason === "clickaway") {
      return;
    }

    setSnackbarOpen(false);
  };

  const copyToClipboard = useCallback(async () => {
    const url = window.location.origin + window.location.pathname;
    await navigator.clipboard.writeText(url);
    setSnackbarOpen(true);
  }, []);

  if (shweetId === undefined || shweet === undefined) {
    // TODO loading
    return <></>;
  }

  return (
    <>
      <Box sx={{ flexGrow: 1 }}>
        <AppBar position="static">
          <Toolbar>
            <IconButton
              edge="start"
              color="inherit"
              aria-label="go back"
              sx={{ mr: 2 }}
              onClick={() => navigate(-1)}
            >
              <ArrowBackIcon />
            </IconButton>

            <Typography variant="h6">Shweet</Typography>
          </Toolbar>
        </AppBar>
      </Box>

      <Container maxWidth="sm" sx={{ padding: "0" }}>
        <Box mb={3} />
      </Container>

      <Container component="main" maxWidth="sm">
        <Box
          sx={{
            marginTop: 2,
            display: "flex",
            flexDirection: "column",
            alignItems: "stretch",
          }}
        >
          {error !== undefined ? (
            <Typography textAlign="center">{error}</Typography>
          ) : (
            <></>
          )}

          <Box display="flex" flexDirection="row">
            <Box mr={2}>
              <UserAvatar user={shweet.user} size={48} />
            </Box>

            <Box
              display="flex"
              flexDirection="column"
              alignItems="stretch"
              flexGrow={1}
            >
              <Box display="flex">
                <Box display="flex" flexGrow={1} flexDirection="column">
                  <Typography>{shweet.user.name}</Typography>
                  <Typography variant="caption">
                    @{shweet.user.username}
                  </Typography>
                </Box>
                <IconButton size="small" onClick={copyToClipboard}>
                  <ShareIcon fontSize="inherit" />
                </IconButton>
              </Box>
            </Box>
          </Box>
          <Box mt={2} mb={1}>
            <Typography variant="body1">{shweet.message}</Typography>
          </Box>

          <Typography variant="caption" justifySelf="end">
            {dayjs(shweet.created_at).format("HH:mm:ss · DD MMM YYYY")}
          </Typography>
        </Box>
        <Box my={1}>
          <Divider />
        </Box>

        {/* TODO:
        [ ] List users who liked/reshweeted
        [ ] Add responses (with mentions)
        */}
        
        <Box display="flex" flexGrow={1} justifyContent="space-evenly">
          <Typography>
            {shweet.reshweet_count}{" "}
            <Typography variant="caption">reshweets</Typography>
          </Typography>
          <Box ml={2} />
          <Typography>
            {shweet.like_count} <Typography variant="caption">likes</Typography>
          </Typography>
          <Box ml={2} />
          <Typography>
            {shweet.comment_count}{" "}
            <Typography variant="caption">comments</Typography>
          </Typography>
        </Box>

        <Box my={1}>
          <Divider />
        </Box>

        <Box display="flex" flexGrow={1} justifyContent="space-evenly">
          <Tooltip title="reshweet">
            <IconButton>
              <CachedIcon />
            </IconButton>
          </Tooltip>
          <Box ml={2} />
          <Tooltip title="like">
            <IconButton onClick={() => likeOrUnlike()}>
              {shweet.liked ? (
                <FavoriteIcon color="error" />
              ) : (
                <FavoriteBorderIcon />
              )}
            </IconButton>
          </Tooltip>
          <Box ml={2} />
          <Tooltip title="comment">
            <IconButton>
              <ChatBubbleOutlineIcon />
            </IconButton>
          </Tooltip>
        </Box>
      </Container>
      <Snackbar
        open={snackbarOpen}
        autoHideDuration={6000}
        onClose={handleClose}
        message="Shweet link copied to your clipboard."
        action={
          <IconButton
            size="small"
            aria-label="close"
            color="inherit"
            onClick={handleClose}
          >
            <CloseIcon fontSize="small" />
          </IconButton>
        }
      />
    </>
  );
};

export default ShweetDetailsPage;

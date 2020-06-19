import React, { useEffect, useState } from "react";
import axios from "axios";
import ListItem from "@material-ui/core/ListItem";
import ListItemText from "@material-ui/core/ListItemText";
import Divider from "@material-ui/core/Divider";
import Typography from "@material-ui/core/Typography";
import Paper from "@material-ui/core/Paper";
import List from "@material-ui/core/List";
import CircularProgress from "@material-ui/core/CircularProgress";

function ListItemLink(props) {
  return <ListItem button component="a" {...props} />;
}

export default function Deals() {
  const [isLoading, setLoading] = useState(true);
  const [state, setState] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      await axios
        .get("/api/v1/deals")
        .then((response) => {
          const data = [];
          response.data.forEach((item) => {
            data.push(item);
          });
          data.sort((a, b) => (a.votes < b.votes ? 1 : -1));
          data.length = 50;
          setState(data);
          setLoading(false);
        })
        .catch(() => {});
    };
    fetchData();
  }, []);

  const getData = state.map((item) => {
    return (
      <div key={item.id}>
        <ListItemLink
          href={item.link}
          target="_blank"
          rel="noopener noreferrer"
        >
          <ListItemText
            disableTypography
            primary={
              <Typography type="body2" style={{ fontWeight: 500 }}>
                {item.title}
              </Typography>
            }
            secondary={`+${item.votes}`}
          />
        </ListItemLink>
        <Divider variant="middle" />
      </div>
    );
  });

  return (
    <>
      <Typography variant="h4" component="h4" gutterBottom>
        Top 50 deals from the last 48 hours
      </Typography>
      <Paper elevation={24}>
        <List dense>
          {isLoading ? (
            <ListItem>
              <CircularProgress />
            </ListItem>
          ) : (
            <>{getData}</>
          )}
        </List>
      </Paper>
    </>
  );
}

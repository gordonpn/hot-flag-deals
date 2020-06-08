import { makeStyles } from "@material-ui/core/styles";

const useStyles = makeStyles({
  spacer: {
    height: "2vh",
  },
});

export default function Spacer() {
  const classes = useStyles();
  return <div className={classes.spacer} />;
}

import React from 'react';

import Card, { CardHeader, CardMedia, CardContent, CardActions } from 'material-ui/Card';
import Button from 'material-ui/Button';

export default class Login extends React.Component {
  render() {
    return (
      <Card style={{ maxWidth: 400, }}>
        <CardHeader title="Login" />
        <CardContent>

        </CardContent>
        <CardActions>
          <Button
            style={{ margin: '0 5px' }}
            color="primary"
            size="small"
            onClick={() => {
              localStorage.setItem('logged', 'true')
              this.props.history.push('/')
            }}> Login</Button>
        </CardActions>
      </Card>
    )
  }
}
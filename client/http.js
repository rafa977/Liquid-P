const express = require('express')
const app = express()
const Api = require('./api/api');
const port = 3000
const cors = require('cors');

app.use(cors());

const winston = require('winston');
const { format } = winston;
const { combine, label, json, splat, timestamp, printf  } = format;

const plainFormat = printf( ({ level, message, timestamp ,metadata}) => {
  let msg = `${timestamp} [${level}] : ${message} `
  if(metadata) {
	   msg += JSON.stringify(metadata)
  }
  return msg
});

const logger = winston.createLogger({
  format: winston.format.combine(
            format.timestamp(),
            plainFormat
        ),
  transports: [
    new (winston.transports.Console)({
      level: 'debug',
      colorize: true,
      handleExceptions: true,
      label: "http"
    }),
    new (winston.transports.File)({
      level: 'debug',
      colorize: false,
      handleExceptions: true,
      label: "http",
      filename: "./logs/liquid-http.log",
      json: false
    })
  ]
});

function getUsernameFromHttpHeader(req, res) {
  var username = req.header("user");
  if (username) {
    return username;
  } else {
    res.status(400).send("missing fabric username in HTTP header 'user'");
  }
}

app.get('/', (req, res) => {
  res.send('The Liquid Gateway Application is running.')
})

app.use(express.json())
app.use('/users', require('./http_user'))

app.post('/postRequestWithJsonData', function(req, res) {
  res.set('Access-Control-Allow-Origin', '*');
    let value

    let requestType = req.body.reqType;
    let username = req.body.username;
    let id = ""

    if(req.body.data.id){
     id = req.body.data.id
    }
    let data = req.body.args

    var allData = data.split(",");

    if(req.body.data.singleInput){
      if(req.body.data.isJson == "yes" && req.body.data.singleInput == "no"){
        value = JSON.stringify(req.body.data)
        allData = new Array(id, value)
      } else if(req.body.data.isJson == "yes" && req.body.data.singleInput == "yes"){
        value = JSON.stringify(req.body.data)
        allData = new Array(value)
      }
    }else{
      if(req.body.data.isJson == "yes"){
        value = JSON.stringify(req.body.data)
        allData = new Array(id, value)
      }
    }

    var startTime, endTime;
    //store the start datetime
    startTime = new Date();
    // var username = ""
    
    // try{
    //   let tokenHeaderKey = process.env.TOKEN_HEADER_KEY;
    //   let jwtSecretKey = process.env.JWT_SECRET_KEY;
    //   const token = req.header(tokenHeaderKey);
    //   username = new Api().validateToken(token, jwtSecretKey)

    // } catch (error) {
    //     logger.error('HTTP/Post ' + requestType + ': A general error occurred:'+ error.message.toString());
    //     endTime = new Date();
    //     res.status(500).json({status: 'error',
    //                             response: null,
    //                             error: error.message.toString(),
    //                             start: startTime.toString(),
    //                             end:endTime.toString(),
    //                             ms_ellapsed: endTime -  startTime
    //                           });
    //     return;
    // }


    try {

        let reqArgs = [];
        reqArgs.push(requestType);
        reqArgs.push(username);
        reqArgs.push(allData);
        console.log(reqArgs)

        new Api().postRequestWithJsonData(reqArgs)
          .then((response) => {
            logger.info('HTTP/Post ' + requestType + ': Callback response'+ JSON.stringify(response.toString()));
            endTime = new Date();
            res.status(200).json({status: 'success',
                                  response: response.toString(),
                                  error: null,
                                  start: startTime.toString(),
                                  end:endTime.toString(),
                                  ms_ellapsed: endTime -  startTime
                                  });
          })
          .catch((error) => {
            logger.error('HTTP/Post ' + requestType + ': An error occurred:'+ error.message.toString());
            endTime = new Date();
            res.status(200).json({status: 'error',
                                  response: null,
                                  error: error.message.toString(),
                                  start: startTime.toString(),
                                  end:endTime.toString(),
                                  ms_ellapsed: endTime -  startTime
                                });
          });

        } catch (error) {
        logger.error('HTTP/Post ' + requestType + ': A general error occurred:'+ error.message.toString());
        endTime = new Date();
        res.status(500).json({status: 'error',
                                response: null,
                                error: error.message.toString(),
                                start: startTime.toString(),
                                end:endTime.toString(),
                                ms_ellapsed: endTime -  startTime
                              });
        return;
      }
  });


  app.post('/getRequestWithArrayArg', function(req, res) {
    res.set('Access-Control-Allow-Origin', '*');
    
    let value

    let requestType = req.body.reqType;
    let username = req.body.username;
    
    // let id = req.body.data.id
    let data = req.body.args

    var allData = data.split(",");
    var startTime, endTime;
    //store the start datetime
    startTime = new Date();
    // var username = ""
    
    // try{
    //   let tokenHeaderKey = process.env.TOKEN_HEADER_KEY;
    //   let jwtSecretKey = process.env.JWT_SECRET_KEY;
    //   const token = req.header(tokenHeaderKey);
    //   username = new Api().validateToken(token, jwtSecretKey)

    // } catch (error) {
    //     logger.error('HTTP/Post ' + requestType + ': A general error occurred:'+ error.message.toString());
    //     endTime = new Date();
    //     res.status(500).json({status: 'error',
    //                             response: null,
    //                             error: error.message.toString(),
    //                             start: startTime.toString(),
    //                             end:endTime.toString(),
    //                             ms_ellapsed: endTime -  startTime
    //                           });
    //     return;
    // }


    try {

        let reqArgs = [];
        reqArgs.push(requestType);
        reqArgs.push(username);
        reqArgs.push(allData);

        new Api().getQueryArray(reqArgs)
          .then((response) => {
            logger.info('HTTP/Get ' + requestType + ': Callback response'+ JSON.stringify(response.toString()));
            endTime = new Date();
            res.status(200).json({status: 'success',
                                  response: response.toString(),
                                  error: null,
                                  start: startTime.toString(),
                                  end:endTime.toString(),
                                  ms_ellapsed: endTime -  startTime
                                  });
          })
          .catch((error) => {
            logger.error('HTTP/Get ' + requestType + ': An error occurred:'+ error.message.toString());
            endTime = new Date();
            res.status(200).json({status: 'error',
                                  response: null,
                                  error: error.message.toString(),
                                  start: startTime.toString(),
                                  end:endTime.toString(),
                                  ms_ellapsed: endTime -  startTime
                                });
          });

        } catch (error) {
        logger.error('HTTP/Get ' + requestType + ': A general error occurred:'+ error.message.toString());
        endTime = new Date();
        res.status(500).json({status: 'error',
                                response: null,
                                error: error.message.toString(),
                                start: startTime.toString(),
                                end:endTime.toString(),
                                ms_ellapsed: endTime -  startTime
                              });
        return;
      }
  });  

  app.get('/getRequestWithArg', function(req, res) {
    res.set('Access-Control-Allow-Origin', '*');
    let requestType = req.query.reqType;
    let id = req.query.id
    let username = req.query.username;

    var startTime, endTime;
    //store the start datetime
    startTime = new Date();
    // var username = ""
    
    // try{
    //   let tokenHeaderKey = process.env.TOKEN_HEADER_KEY;
    //   let jwtSecretKey = process.env.JWT_SECRET_KEY;
    //   const token = req.header(tokenHeaderKey);
    //   username = new Api().validateToken(token, jwtSecretKey)

    // } catch (error) {
    //     logger.error('HTTP/Post ' + requestType + ': A general error occurred:'+ error.message.toString());
    //     endTime = new Date();
    //     res.status(500).json({status: 'error',
    //                             response: null,
    //                             error: error.message.toString(),
    //                             start: startTime.toString(),
    //                             end:endTime.toString(),
    //                             ms_ellapsed: endTime -  startTime
    //                           });
    //     return;
    // }

    try {

        let reqArgs = [];
        reqArgs.push(username);
        reqArgs.push(id);
        reqArgs.push(requestType);

        new Api().getQuery(reqArgs)
          .then((response) => {
            logger.info('HTTP/Get ' + requestType + ': Callback response'+ JSON.stringify(response.toString()));
            endTime = new Date();
            res.status(200).json({status: 'success',
                                  response: response.toString(),
                                  error: null,
                                  start: startTime.toString(),
                                  end:endTime.toString(),
                                  ms_ellapsed: endTime -  startTime
                                  });
          })
          .catch((error) => {
            logger.error('HTTP/Get ' + requestType + ': An error occurred:'+ error.message.toString());
            endTime = new Date();
            res.status(200).json({status: 'error',
                                  response: null,
                                  error: error.message.toString(),
                                  start: startTime.toString(),
                                  end:endTime.toString(),
                                  ms_ellapsed: endTime -  startTime
                                });
          });

        } catch (error) {
        logger.error('HTTP/Get ' + requestType + ': A general error occurred:'+ error.message.toString());
        endTime = new Date();
        res.status(500).json({status: 'error',
                                response: null,
                                error: error.message.toString(),
                                start: startTime.toString(),
                                end:endTime.toString(),
                                ms_ellapsed: endTime -  startTime
                              });
        return;
      }
  });      


  app.get('/verifyDocument', function(req, res) {
    res.set('Access-Control-Allow-Origin', '*');
    let username = req.query.username;
    let docHash = req.query.hash
    let owner = req.query.owner
    let docId = req.query.docId

    try {
        var startTime, endTime;
        //store the start datetime
        startTime = new Date();

        let reqArgs = [];
        reqArgs.push(username);
        reqArgs.push(docHash);
        reqArgs.push(owner);
        reqArgs.push(docId);

        new Api().verifyDocument(reqArgs)
          .then((response) => {
            logger.info('HTTP/Verify Document: Callback response'+ JSON.stringify(response.toString()));
            endTime = new Date();
            res.status(200).json({status: 'success',
                                  response: response.toString(),
                                  error: null,
                                  start: startTime.toString(),
                                  end:endTime.toString(),
                                  ms_ellapsed: endTime -  startTime
                                  });
          })
          .catch((error) => {
            logger.error('HTTP/Verify Document: An error occurred:'+ error.message.toString());
            endTime = new Date();
            res.status(200).json({status: 'error',
                                  response: null,
                                  error: error.message.toString(),
                                  start: startTime.toString(),
                                  end:endTime.toString(),
                                  ms_ellapsed: endTime -  startTime
                                });
          });

        } catch (error) {
        logger.error('HTTP/Verify Document: A general error occurred:'+ error.message.toString());
        endTime = new Date();
        res.status(500).json({status: 'error',
                                response: null,
                                error: error.message.toString(),
                                start: startTime.toString(),
                                end:endTime.toString(),
                                ms_ellapsed: endTime -  startTime
                              });
        return;
      }
  });   
  

app.listen(port, () => {
  console.log(`HTTP: The Liquid Gateway Application is listening at http://localhost:${port}`)
});

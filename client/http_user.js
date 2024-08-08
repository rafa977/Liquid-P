const express = require('express')
var router = require('express').Router();

const app = express()
const Api = require('./api/api');
const port = 3000
const cors = require('cors');

app.use(cors());

const winston = require('winston');
const { format } = winston;
const { combine, label, json, splat, timestamp, printf  } = format;
const dotenv = require('dotenv');
const jwt = require('jsonwebtoken');

dotenv.config();

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
      filename: "./logs/liquid-http_user.log",
      json: false
    })
  ]
});


router.post("/generateToken", function(req, res) {
  res.set('Access-Control-Allow-Origin', '*');
  // Validate User Here
  // Then generate JWT Token

  let jwtSecretKey = process.env.JWT_SECRET_KEY;
  let data = {
      username: "Kat"
  }

  const token = jwt.sign(data, jwtSecretKey);

  res.send(token);
});




router.post('/registerUser', function(req, res) {
    res.set('Access-Control-Allow-Origin', '*');

    let value
    var startTime, endTime;

    startTime = new Date();

    let id = req.body.data.id;
    let isSuperAdmin = req.body.isSuperAdmin;
    
    req.body.data.date = startTime;

    value = JSON.stringify(req.body.data)
    allData = new Array(id, value)
    
    let username = req.body.username;
    let password = req.body.password;
    let ca = req.body.ca;
    let caAdmin = req.body.caAdmin;
    let mspId = req.body.mspId;

    try {
        var requestType = "checkIdExists";

        var reqArgs = [];
        reqArgs.push("Admin");
        reqArgs.push(id);
        reqArgs.push(requestType);

        new Api().getQuery(reqArgs)
          .then((response) => {
            logger.info('HTTP/Get ' + requestType + ': Callback response'+ JSON.stringify(response.toString()));         

            reqArgs = [];
            reqArgs.push(username);
            reqArgs.push(password);
            reqArgs.push(isSuperAdmin);
            reqArgs.push(ca);
            reqArgs.push(caAdmin);
            reqArgs.push(mspId);

            new Api().registerUser(reqArgs)
            .then((response) => {
              logger.info('HTTP/Post : Callback response'+ JSON.stringify(response.toString()));
              endTime = new Date();
  
              let registerArgs = [];
              registerArgs.push("registerAccount");
              registerArgs.push(username);
              registerArgs.push(allData);
              

              requestType = "registerAccount";
  
              new Api().postRequestWithJsonData(registerArgs)
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
  
                  new Api().removeFromWallet(username)
              });
            })
            .catch((error) => {
              logger.error('HTTP/Post: An error occurred:'+ error.message.toString());
              endTime = new Date();
  
              res.status(200).json({status: 'error',
                                    response: null,
                                    error: error.message.toString(),
                                    start: startTime.toString(),
                                    end:endTime.toString(),
                                    ms_ellapsed: endTime -  startTime
                                  });
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
        logger.error('HTTP/Post : A general error occurred:'+ error.message.toString());
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



  router.post('/updateUserAttr', function(req, res) {
    res.set('Access-Control-Allow-Origin', '*');

    var startTime, endTime;

    startTime = new Date();

    let isSuperAdmin = req.body.isSuperAdmin;
    let caller = req.body.caller;

    let username = req.body.username;
    let password = req.body.password;
    let ca = req.body.ca;
    let caAdmin = req.body.caAdmin;
    let mspId = req.body.mspId;

    try {
        var requestType = "udpateAttr";

        var reqArgs = [];
        reqArgs.push(caller);
        reqArgs.push(username);
        reqArgs.push(password);
        reqArgs.push(isSuperAdmin);
        reqArgs.push(ca);
        reqArgs.push(caAdmin);
        reqArgs.push(mspId); 

        new Api().updateUserAttribute(reqArgs)
          .then((response) => {
            logger.info('HTTP/Get ' + requestType + ': Callback response'+ JSON.stringify(response.toString()));         

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
        logger.error('HTTP/Post : A general error occurred:'+ error.message.toString());
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


  router.post('/removeUser', function(req, res) {
    res.set('Access-Control-Allow-Origin', '*');

    let caller = req.body.caller;
    let username = req.body.username;
    let caName = req.body.caName;

    let data = req.body.args
    var allData = data.split(",");
    

    try {
        var startTime, endTime;
        //store the start datetime
        startTime = new Date();

        let registerArgs = [];
        registerArgs.push("permDeleteAccount");
        registerArgs.push(caller);
        registerArgs.push(allData);

        var requestType = "permDeleteAccount";

        // //Permanent delete from Blockchain user data
        new Api().postRequestWithJsonData(registerArgs)
          .then((response) => {
            logger.info('HTTP/Post : Callback response'+ JSON.stringify(response.toString()));

            endTime = new Date();

            //Remove & revoke user from Fabric and Wallet
            new Api().removeUser(username, caName)
            .then(() => {

              endTime = new Date();
              res.status(200).json({status: 'success',
                                    response: null,
                                    error: null,
                                    start: startTime.toString(),
                                    end:endTime.toString(),
                                    ms_ellapsed: endTime -  startTime
                                    });
            })
            .catch((error) => {

              endTime = new Date();
              res.status(200).json({status: 'error',
                                    response: null,
                                    error: error.message.toString(),
                                    start: startTime.toString(),
                                    end:endTime.toString(),
                                    ms_ellapsed: endTime -  startTime
                                  });
            });
          })
          .catch((error) => {
            logger.error('HTTP/Post: An error occurred:'+ error.message.toString());
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
        logger.error('HTTP/Post : A general error occurred:'+ error.message.toString());
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

  router.get('/revokeUser', function(req, res) {
    res.set('Access-Control-Allow-Origin', '*');
    let value

    let username = req.query.username;

    try {
        var startTime, endTime;
        //store the start datetime
        startTime = new Date();

        new Api().revokerUser(username)
          .then((response) => {
            logger.info('HTTP/Post : Callback response'+ JSON.stringify(response.toString()));
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
            logger.error('HTTP/Post: An error occurred:'+ error.message.toString());
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
        logger.error('HTTP/Post : A general error occurred:'+ error.message.toString());
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


  router.get('/enrollAdmin', function(req, res) {
    res.set('Access-Control-Allow-Origin', '*');
    let value

    let ca = req.query.ca;
    let adminUsername = req.query.adminUsername;
    let mspId = req.query.mspId;
    

    try {
        var startTime, endTime;
        //store the start datetime
        startTime = new Date();

        let reqArgs = [];
        reqArgs.push(ca);
        reqArgs.push(adminUsername);
        reqArgs.push(mspId);

        new Api().enrollAdmin(reqArgs)
          .then((response) => {
            logger.info('HTTP/Post : Callback response'+ JSON.stringify(response.toString()));
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
            logger.error('HTTP/Post: An error occurred:'+ error.message.toString());
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
        logger.error('HTTP/Post : A general error occurred:'+ error.message.toString());
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

router.get('/removeUserFromFabric', function(req, res) {
    res.set('Access-Control-Allow-Origin', '*');

    let username = req.query.username;    

    try {
        var startTime, endTime;
        //store the start datetime
        startTime = new Date();

        //Permanent delete from Blockchain user data
        new Api().removeUser(username)
          .then((response) => {

            endTime = new Date();
            res.status(200).json({status: 'success',
                response: "Successfully Revoked",
                error: null,
                start: startTime.toString(),
                end:endTime.toString(),
                ms_ellapsed: endTime -  startTime
            });  
        })
          .catch((error) => {
            logger.error('HTTP/Post: An error occurred here:'+ error.message.toString());
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
        logger.error('HTTP/Post : A general error occurred:'+ error.message.toString());
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

module.exports = router;

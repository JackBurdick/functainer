#!/usr/bin/env python3


import http.server
import socketserver

class HttpHandler(http.server.SimpleHTTPRequestHandler):

    def do_GET(self):
        if self.path == "/tf_cpu_classification":
            r_headers = str(self.headers)
            #data_string = self.rfile.read(int(self.headers['Content-Length']))

            self.send_response(200)
            self.send_header("Content-type", "text/html")
            self.end_headers()
            self.wfile.write(bytes("<html> <head><title> Testing </title> </head><body><p> some sweet classification will happen here.... </p>"+r_headers+"</body>", "UTF-8"))
        else:
            http.server.SimpleHTTPRequestHandler.do_GET(self)

http_server=socketserver.TCPServer(("",8080), HttpHandler)
http_server.serve_forever()
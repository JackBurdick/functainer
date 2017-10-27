#!/usr/bin/env python3

import http.server
import socketserver

class HttpHandler(http.server.SimpleHTTPRequestHandler):

    def do_GET(self):
        if self.path == "/tf_cpu_classification":
            self.send_response(200)
            self.send_header("Content-type", "text/html")
            self.end_headers()
            self.wfile.write(bytes("<html> <head><title> Testing </title> </head><body><p> some sweet classification will happen here.... </p>"+str(self.headers)+"</body>", "UTF-8"))
        else:
            http.server.SimpleHTTPRequestHandler.do_GET(self)

http_server=socketserver.TCPServer(("",8080), HttpHandler)
http_server.serve_forever()
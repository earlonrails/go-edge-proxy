#!/usr/bin/env ruby
require 'redis'
require 'json'

sub_client = Redis.new
pub_client = Redis.new
count = 0
sub_client.psubscribe('pre-filter:*') do |on|
  on.pmessage do |pattern, channel, request_id|
    puts "received message: #{request_id}"
    pub_client.setex(request_id, 1, {foo: "bar", baz: request_id}.to_json)
    puts "set complete"
    pub_client.publish("respond:#{request_id}", request_id)
    puts "publish to respond"
    puts "#{count += 1}"
  end
end

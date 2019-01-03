#!/usr/bin/env ruby
require 'redis'
require 'json'
require 'logger'

sub_client = Redis.new(url: 'redis://redis:6379')
pub_client = Redis.new(url: 'redis://redis:6379')
count = 0
logger = Logger.new('/proc/1/fd/1')
logger.info "Worker started"
sub_client.psubscribe('pre-filter:*') do |on|
  on.pmessage do |pattern, channel, request_id|
    logger.info "received message: #{request_id}"
    pub_client.setex(request_id, 1, {foo: "bar", baz: request_id}.to_json)
    logger.info "set complete"
    pub_client.publish("respond:#{request_id}", request_id)
    logger.info "publish to respond"
    logger.info "#{count += 1}"
  end
end

FROM ruby:2.5.1
RUN mkdir /app
WORKDIR /app
COPY Gemfile /app/Gemfile
COPY Gemfile.lock /app/Gemfile.lock
RUN bundle install
RUN bundle exec rails new demo
WORKDIR /app/demo
CMD bundle exec rails s -p 3000 -b '0.0.0.0'

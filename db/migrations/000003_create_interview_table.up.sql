create table if not exists interview (
     id UUID primary key not null,
     user_id UUID not null,
     status text not null,
     job_position text default null,
     job_level text default null,
     question_count integer default 0,

     foreign key (user_id) references "user" (id)
)
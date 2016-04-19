sched
=====

Package sched is used to find the next time event from multiple rules.

A Rule is a simple and flexible interface with one method:

	Next(now time.Time) time.Time

Schedule is a priority queue of rules and their next occurance.
Scheduler is thread safe service to manage a schedule and trigger on events.

	s := sched.New(time.Now(), func(e sched.Event){
		fmt.Println("happy birthday!", e.Time)
	})
	s.Add(myBirthday)

License
-------
sched is BSD licensed, Copyright (c) 2016 Martin Schnabel
